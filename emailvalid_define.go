package emailvalid

import "sync"

type EmailValid struct {
	OriginEmail string   // "Joe Lamèche" <joe.lameche@test.fr>
	Email       string   // joe.lameche@test.fr
	Name        string   // Joe Lamèche
	Local       string   // joe.lameche
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
	TIMEOUT_SPF   uint = 5
	TIMEOUT_DMARC uint = 5
	TIMEOUT_MX    uint = 5
)
