package emailvalid

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"slices"
	"strings"
	"time"
)

// Add ajoute deux entiers et retourne leur somme.
//
// Parameters:
//
//	a: le premier entier à ajouter.
//	b: le deuxième entier à ajouter.
//
// Returns:
//
//	La somme de a et b.
func NewTLD() (tld *Tld, err error) {

	defer func() {
		if errRecover := recover(); errRecover != nil {
			tld = nil
			err = fmt.Errorf("%s", errRecover)
		}
	}()

	if len(tld_default_data) == 0 {
		return nil, fmt.Errorf("no default TLD available")
	}

	return &Tld{
		src:  tld_default.src,
		date: tld_default.date,
		tld:  &tld_default_data,
	}, nil

}

func NewTLDFromFile(fileName string) (tld *Tld, err error) {

	defer func() {
		if errRecover := recover(); errRecover != nil {
			tld = nil
			err = fmt.Errorf("%s", errRecover)
		}
	}()

	data, err := importTLDFromFile(fileName)
	if err != nil {
		return nil, err
	}

	return &Tld{
		src:  strings.TrimSpace(fileName),
		date: time.Now(),
		tld:  data,
	}, nil

}

func NewTLDFromSlice(slice []string) (tld *Tld, err error) {

	defer func() {
		if errRecover := recover(); errRecover != nil {
			tld = nil
			err = fmt.Errorf("%s", errRecover)
		}
	}()

	data, err := importTLDFromSlice(slice)
	if err != nil {
		return nil, err
	}

	return &Tld{
		src:  "SLICE",
		date: time.Now(),
		tld:  data,
	}, nil

}

//

func NewTLDFromUrl[Options string | map[string]any](opt ...Options) (tld *Tld, err error) {

	defer func() {
		if errRecover := recover(); errRecover != nil {
			tld = nil
			err = fmt.Errorf("%s", errRecover)
		}
	}()

	// Init des valeurs par défaut
	var url string
	var timeOut uint16 = IMPORT_TLD_TIMEOUT
	var sizeMax uint32 = IMPORT_TLD_SIZE_MAX
	var tlsSkipVerify bool = IMPORT_TLD_TLS_SKIP_VERIFY
	var verifyContentLength bool = IMPORT_VERIFY_CONTENT_LENGTH

	if len(opt) > 1 {
		return nil, fmt.Errorf("invalid number of arg")
	}

	if len(opt) == 1 {

		switch optValue := any(opt).(type) {

		case []string: // Par défaut on demande l'URL
			url = strings.TrimSpace(optValue[0])

		case []map[string]any:

			for kMap, vMap := range optValue[0] {

				switch kMap {

				case "url":

					if reflect.TypeOf(vMap) != reflect.TypeOf("") {
						return nil, fmt.Errorf("key %s, string variable type expected", kMap)
					}

					v := strings.TrimSpace(vMap.(string))

					if len(v) == 0 {
						return nil, fmt.Errorf("key %s cannot be empty", kMap)
					}
					url = v

				case "timeout":

					if reflect.TypeOf(vMap) != reflect.TypeOf(int(0)) {
						return nil, fmt.Errorf("key %s, integer variable type expected", kMap)
					}

					if vMap.(int) < 0 || vMap.(int) > 3600 {
						return nil, fmt.Errorf("key %s, value must be between [0-3600]", kMap)
					}

					timeOut = uint16(vMap.(int))

				case "sizeMax":

					if reflect.TypeOf(vMap) != reflect.TypeOf(int(0)) {
						return nil, fmt.Errorf("key %s, integer variable type expected", kMap)
					}

					if vMap.(int) <= 0 || vMap.(int) > 1000000 {
						return nil, fmt.Errorf("key %s, value must be between [1-1000000]", kMap)
					}

					sizeMax = uint32(vMap.(int))

				case "SelfCertificate":

					if reflect.TypeOf(vMap) != reflect.TypeOf(bool(false)) {
						return nil, fmt.Errorf("boolean variable type expected")
					}

					tlsSkipVerify = vMap.(bool)

				case "VerifyContentLength":

					if reflect.TypeOf(vMap) != reflect.TypeOf(bool(false)) {
						return nil, fmt.Errorf("boolean variable type expected")
					}

					verifyContentLength = vMap.(bool)

				default:
					return nil, fmt.Errorf("key %s not exist", kMap)
				}

			}

		default:
			return nil, fmt.Errorf("invalid arg")

		}

	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: tlsSkipVerify},
	}

	httpClient := &http.Client{
		Timeout:   time.Duration(timeOut) * time.Second,
		Transport: transport,
	}

	// Obtenir la taille du fichier avant de télécharger
	resp, err := httpClient.Head(url)
	if err != nil {
		return nil, err
	}

	// Check server réponse
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status for url[%s] - %s", url, resp.Status)
	}

	// On vérifie la valeur de ContentLength
	if verifyContentLength {

		// -1: taille inconnue
		if resp.ContentLength < 0 {
			return nil, fmt.Errorf("ContentLength of request GET[%s] length is unknown", url)
		}

		if resp.ContentLength > int64(sizeMax) {
			return nil, fmt.Errorf("max allowed size is %d - URL size is %d", sizeMax, resp.ContentLength)
		}

	}

	if !strings.Contains(strings.ToLower(resp.Header.Get("Content-Type")), IMPORT_TLD_MIME_TYPE) {
		return nil, fmt.Errorf("mimetype invalid")
	}

	resp, err = httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check server réponse
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status for url[%s] - %s", url, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	allTld := make(map[string]byte)

	// Créer un scanner pour lire les données ligne par ligne
	fScan := bufio.NewScanner(bytes.NewReader(data))
	fScan.Split(bufio.ScanLines)

	for fScan.Scan() {

		// By Pass espaces et commentaires
		line := strings.TrimSpace(fScan.Text())
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		tld, err := tldValid(line)
		if err != nil {
			return nil, err
		}

		allTld[tld] = 1
	}

	if err := fScan.Err(); err != nil {
		return nil, err // Gérer les erreurs de scanner
	}

	if len(allTld) == 0 {
		return nil, fmt.Errorf("no TLD")
	}

	return &Tld{
		src:  url,
		date: time.Now(),
		tld:  &allTld,
	}, nil

}

// Compter le nombre de TLD
// Parameters:
//
// Returns:
// Le nombre de TLD existants
// -1: si erreur
func (t *Tld) Count() int {

	defer func() {
		if errRecover := recover(); errRecover != nil {
			_ = -1
		}
	}()

	t.mu.Lock()
	defer t.mu.Unlock()

	return len((*t.tld))
}

// Obtenir la liste des TLD non triés
func (t *Tld) GetAllTLD() (tld []string, err error) {

	defer func() {
		if errRecover := recover(); errRecover != nil {
			tld = nil
			err = fmt.Errorf("%s", errRecover)
		}
	}()

	// Pré-allocation de la liste
	tld = make([]string, 0, len(*t.tld))

	t.mu.Lock()
	defer t.mu.Unlock()

	for k := range *t.tld {
		tld = append(tld, k)
	}

	return tld, nil

}

// Obtenir la liste des TLD triés par ordre croissant
func (t *Tld) GetAllTLDInc() ([]string, error) {

	list, err := t.GetAllTLD()
	if err != nil {
		return nil, err
	}

	slices.Sort(list)

	return list, nil
}

// Obtenir la liste des TLD triés par ordre décroissant
func (t *Tld) GetAllTLDDec() ([]string, error) {

	list, err := t.GetAllTLDInc()
	if err != nil {
		return nil, err
	}

	slices.Reverse(list)

	return list, nil
}

// Vérifier la présence d'un TLD
func (t *Tld) IsTLD(data string) bool {

	defer func() {
		if errRecover := recover(); errRecover != nil {
			_ = false
		}
	}()

	t.mu.Lock()
	defer t.mu.Unlock()

	return (*t.tld)[strings.ToUpper(data)] == 1
}

// Vérifier la Syntaxe d'un supposé TLD
// OK -> TLD retourné et err == nil / KO -> nil, error
func CheckSyntaxeTLD(tld string) (string, error) {
	return tldValid(tld)

}

// Ajouter un TLD
// Si erreur de syntaxe retourne une erreur.
func (t *Tld) Add(tld string) error {

	defer func() {
		if errRecover := recover(); errRecover != nil {
			_ = fmt.Errorf("%s", errRecover)
		}
	}()

	data, err := tldValid(strings.TrimSpace(tld))
	if err != nil {
		return err
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	(*t.tld)[strings.ToUpper(data)] = 1

	return nil
}

// Supprimer un TLD
func (t *Tld) Del(tld string) error {

	defer func() {
		if errRecover := recover(); errRecover != nil {
			_ = fmt.Errorf("%s", errRecover)
		}
	}()

	txtUpper, err := tldValid(strings.TrimSpace(tld))
	if err != nil {
		return err
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	_, exist := (*t.tld)[txtUpper]
	if !exist {
		return fmt.Errorf("tld %s not exist", tld)
	}

	delete(*t.tld, txtUpper)

	_, exist = (*t.tld)[txtUpper]
	if exist {
		return fmt.Errorf("tld %s already exist", txtUpper)
	}

	return nil
}
