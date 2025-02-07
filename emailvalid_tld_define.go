package emailvalid

import (
	"sync"
	"time"
)

const (
	IMPORT_TLD_TIMEOUT           uint16 = 3     // secondes [0;3600]
	IMPORT_TLD_SIZE_MAX          uint32 = 10000 // octets MAX: 1000000
	IMPORT_TLD_TLS_SKIP_VERIFY   bool   = false
	IMPORT_VERIFY_CONTENT_LENGTH bool   = false
	IMPORT_TLD_MIME_TYPE         string = "text/plain"
)

type Tld struct {
	mu   sync.Mutex
	src  string
	date time.Time
	tld  *map[string]byte
}
