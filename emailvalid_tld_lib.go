package emailvalid

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func importTLDFromFile(file string) (*map[string]byte, error) {

	file = strings.TrimSpace(file)

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	mimetype, err := getFileContentType(f)
	if err != nil {
		return nil, err
	}

	if !strings.Contains(mimetype, "text/plain") {
		return nil, fmt.Errorf("file '%s' must be text/plain type", file)
	}

	var ligne uint64 = 0

	tld := make(map[string]byte)

	fScan := bufio.NewScanner(f)
	fScan.Split(bufio.ScanLines)

	for fScan.Scan() {

		ligne++

		// By Pass espaces et commentaires
		txt := strings.TrimSpace(fScan.Text())
		if len(txt) == 0 || txt[0] == '#' {
			continue
		}

		_, err := tldValid(txt)
		if err != nil {
			return nil, fmt.Errorf("file: %s - ligne: %d - TLD: %s : %s", file, ligne, txt, err.Error())
		}

		txtUpper := strings.ToUpper(txt)

		_, exist := tld[txtUpper]
		if exist {
			return nil, fmt.Errorf("file: %s - ligne: %d - TLD: %s : already exist", file, ligne, txt)
		}

		tld[txtUpper] = 1
	}

	if err := fScan.Err(); err != nil {
		return nil, err
	}

	if (len(tld)) == 0 {
		return nil, fmt.Errorf("import TLD form slice is zero")
	}

	return &tld, nil
}

func tldValid(data string) (string, error) {

	data = strings.ToUpper(strings.TrimSpace(data))

	// Vérifier la taille du TLD [2-63]
	if len(data) < 2 || len(data) > 63 {
		return "", errors.New("tld length not valid")
	}

	// Punycode
	if strings.HasPrefix(data, "XN--") {

		if !regexp.MustCompile(`^(?:(XN--[A-Z0-9]+(?:-[A-Z0-9]+)*[A-Z0-9]?))$`).MatchString(data) {
			return "", errors.New("tld syntax not valid")
		}

	} else {

		if !regexp.MustCompile(`^(?:([A-Z][A-Z0-9]*(?:-[A-Z0-9]+)*[A-Z0-9]?))$`).MatchString(data) {
			return "", errors.New("tld syntax not valid")
		}

	}

	return data, nil
}

func getFileContentType(file *os.File) (string, error) {

	defer func() {
		// Retour du descripteur au debut du fichier
		file.Seek(0, io.SeekStart)
	}()

	buf := make([]byte, 512)
	_, err := file.Read(buf)
	if err != nil {
		return "", err
	}

	return strings.ToLower(strings.TrimSpace(http.DetectContentType(buf))), nil
}

func importTLDFromSlice(data []string) (*map[string]byte, error) {

	tld := map[string]byte{}

	for k, v := range data {

		vv := strings.ToUpper(strings.TrimSpace(v))

		_, exist := tld[vv]
		if exist {
			return nil, fmt.Errorf("entrie: %d - TLD: %s - already exist", k, v)
		}

		_, err := tldValid(v)
		if err != nil {
			return nil, fmt.Errorf("entrie: %d - TLD: %s - %s", k, v, err.Error())
		}

		tld[vv] = 1

	}

	if (len(tld)) == 0 {
		return nil, fmt.Errorf("import TLD form slice is zero")
	}

	return &tld, nil
}
