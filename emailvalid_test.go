package emailvalid_test

import (
	"fmt"
	myEmailValid "matyas-cyril/email-valid"
	"net"
	"strings"
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

	for dom, mx := range mapMX {

		for host, val := range mx {

			fmt.Println(dom, "-->", host, ":")

			for k, v := range val.(map[string]any) {
				switch strings.ToUpper(k) {

				case "IP":
					fmt.Println("\tIP:", v.([]net.IP))

				case "PREF":
					fmt.Println("\tPref:", v.(uint16))
				}
			}

		}

	}

}
