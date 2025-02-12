package emailvalid_test

import (
	"fmt"
	myEmailValid "matyas-cyril/email-valid"
	"testing"
)

const EMAIL string = "\"Joe Lamèche\" <joe.lameche@test.fr>"

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

	mapMX, err := e.GetMX()
	if err != nil {
		t.Fatal(err)
	}

	for dom, tabMX := range mapMX {

		fmt.Printf("%s:\n", dom)

		for _, mx := range tabMX {

			fmt.Printf("\t%s -> @IP:%v\n", mx.Host, mx.IP)
			fmt.Printf("\t%s -> Pref:%d\n", mx.Host, mx.Pref)

		}

	}

}
