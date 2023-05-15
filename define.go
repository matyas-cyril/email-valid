package emailvalid

import "sync"

type Tld struct {
	mu  sync.Mutex
	tld *map[string]byte
}

type EmailValid struct {
	OriginEmail string   // "Joe Lamoustache" <joe.lamoustache@test.fr>
	Email       string   // joe.lamoustache@test.fr
	Name        string   // Joe Lamoustache
	Local       string   // joe.lamoustache
	Domain      string   // TEST.FR
	DomArrayInv []string // [FR TEST]
	timeout     struct {
		dmarc uint
		mx    uint
		spf   uint
	}
	mu sync.Mutex
}

const (
	regexTLD      string = `^(?:[A-Z]{2,63}|XN\-\-(?:[A-Z0-9]{2,59}|[A-Z][A-Z0-9]+\-[A-Z]{2,56}))$`
	regexMail     string = `^(.+)\@([^\@\s]{1,63})$`
	timeout_spf   uint   = 5
	timeout_dmarc uint   = 5
	timeout_mx    uint   = 5
)
