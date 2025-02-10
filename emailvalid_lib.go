package emailvalid

import (
	"net/mail"
	"regexp"
	"strings"
)

// Convertir un email (RFC5322) en structure EmailValid
func extractEmail(email string) (*EmailValid, error) {

	email = strings.TrimSpace(email)

	e, err := mail.ParseAddress(email)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`^(.+)\@([^\@\s]{1,663})$`)
	m := re.FindStringSubmatch(e.Address)
	if len(m) == 3 {

		domArrayInv := splitReverseDomain(m[2])

		return &EmailValid{
				srcEmail:    email,
				email:       strings.ToLower(strings.TrimSpace(e.Address)),
				name:        strings.TrimSpace(e.Name),
				local:       strings.ToLower(strings.TrimSpace(m[1])),
				domain:      strings.ToUpper(strings.TrimSpace(m[2])),
				domArrayInv: domArrayInv,
				timeout: struct {
					dmarc uint
					mx    uint
					spf   uint
				}{
					dmarc: TIMEOUT_DMARC,
					mx:    TIMEOUT_MX,
					spf:   TIMEOUT_SPF,
				},
			},
			nil

	}

	return nil, nil
}

// splitReverseDomain décompose un FQDN
// et va présenter le TLD en 1ere case
// Ex: test.fr -> [FR TEST]
func splitReverseDomain(domain string) []string {

	domain = strings.ToUpper(strings.TrimSpace(domain))
	if len(domain) == 0 {
		return nil
	}

	var invStringArray = func(array []string) []string {

		var newArray []string
		for i := (len(array) - 1); i >= 0; i = i - 1 {

			if data := strings.TrimSpace(array[i]); len(data) > 0 {
				newArray = append(newArray, data)
			}

		}
		return newArray
	}

	split := invStringArray(strings.Split(domain, "."))
	if string(domain[len(domain)-1]) == "." {
		return split[1:]
	}

	return split
}
