package emailvalid

import "sync"

type EmailValid struct {
	srcEmail    string   // "Joe Lamèche" <joe.lameche@test.fr>
	email       string   // joe.lameche@test.fr
	name        string   // Joe Lamèche
	local       string   // joe.lameche
	domain      string   // TEST.FR
	domArrayInv []string // [FR TEST]
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
