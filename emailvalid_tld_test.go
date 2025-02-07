package emailvalid_test

import (
	"fmt"
	myEmailValid "matyas-cyril/email-valid"
	"testing"
)

var (
	// Liste des données pour les tests.
	dataList = []string{"Ca", "PL", "cnx", "FR", "uK", "beer", "looser"}

	// Nom du fichier d'import pour les tests
	dataFile = "tlds-alpha-by-domain.txt"

	// Adresse d'obtention des TLD
	url string = "https://data.iana.org/TLD/tlds-alpha-by-domain.txt"
)

// go test -timeout 3s -run ^TestNewTLD$
func TestNewTLD(t *testing.T) {

	// Avec les valeurs par défaut
	fmt.Print("Load from Default -> ")
	tld, err := myEmailValid.NewTLD()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("entries : %d\n", tld.Count())

	// Lister par order décroissant
	d, err := tld.GetAllTLDDec()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("GetAllTLDDec : %s", d)

}

// go test -timeout 3s -run ^TestNewTLDFromFile$
func TestNewTLDFromFile(t *testing.T) {

	fmt.Print("Load from File -> ")
	tld, err := myEmailValid.NewTLDFromFile(dataFile)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("entries : %d\n", tld.Count())

}

// go test -timeout 3s -run ^TestNewTLDFromSlice$
func TestNewTLDFromSlice(t *testing.T) {

	fmt.Print("Load from Slice -> ")
	tld, err := myEmailValid.NewTLDFromSlice(dataList)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("entries : %d\n", tld.Count())

}

// go test -timeout 3s -run ^TestIsTLD$
func TestIsTLD(t *testing.T) {

	tld, err := myEmailValid.NewTLD()
	if err != nil {
		t.Fatal(err)
	}

	for _, i := range dataList {
		fmt.Printf("%s is valid TLD : %v\n", i, tld.IsTLD(i))
	}

}

// go test -timeout 3s -run ^TestADDTLD$
func TestADDTLD(t *testing.T) {

	data := "Taupe"

	tld, err := myEmailValid.NewTLD()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("entries : %d\n", tld.Count())
	fmt.Printf("%s is valid TLD : %v\n", data, tld.IsTLD(data))

	err = tld.Add(data)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("entries : %d\n", tld.Count())
	fmt.Printf("%s is valid TLD : %v\n", data, tld.IsTLD(data))
}

// go test -timeout 3s -run ^TestDELTLD$
func TestDELTLD(t *testing.T) {

	data := "beer"

	tld, err := myEmailValid.NewTLD()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("entries : %d\n", tld.Count())
	fmt.Printf("%s is valid TLD : %v\n", data, tld.IsTLD(data))

	err = tld.Del(data)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("entries : %d\n", tld.Count())
	fmt.Printf("%s is valid TLD : %v\n", data, tld.IsTLD(data))
}

// go test -timeout 3s -run ^TestCheckSyntaxeTLD$
func TestCheckSyntaxeTLD(t *testing.T) {

	data := []string{"abc123", "123abc", "tld!", "-tld", "tld-", "a", "123", "example@domain", "ex ample", "Io"}

	for _, i := range data {

		_, err := myEmailValid.CheckSyntaxeTLD(i)
		if err != nil {
			fmt.Printf("KO -> '%s' is not a valid TLD\n", i)
			continue
		}
		fmt.Printf("OK -> '%s' is a valid TLD\n", i)

	}

}

// go test -timeout 3s -run ^TestTLDFromUrl$
func TestTLDFromUrl(t *testing.T) {

	fmt.Printf("URL:%s\n", url+"e")

	// Test avec une URL invalide
	_, err := myEmailValid.NewTLDFromUrl(url + "e")
	if err != nil {
		fmt.Printf("Err: %s\n\n", err)
	}

	tld, err := myEmailValid.NewTLDFromUrl(url)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("entries : %d\n", tld.Count())

	// Lister par order décroissant
	d, err := tld.GetAllTLDDec()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("GetAllTLDDec : %s", d)

}
