package emailvalid_test

import (
	"fmt"
	myEmailValid "matyas-cyril/email-valid"
	"testing"
)

const EMAIL string = "\"Joe Lamèche\" <joe.lameche@free.fr>"

// go test -timeout 3s -run ^TestNewEmail$
func TestNewEmail(t *testing.T) {

	e, err := myEmailValid.NewEmail(EMAIL)
	if err != nil {
		t.Fatal(err)
	}

	d, err := e.GetAllEmailData()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(d)
}

// go test -timeout 3s -run ^TestGetMX$
func TestGetMX(t *testing.T) {

	e, err := myEmailValid.NewEmail(EMAIL)
	if err != nil {
		t.Fatal(err)
	}

	json, err := e.GetMX()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(json))
}
