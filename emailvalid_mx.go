package emailvalid

import (
	"context"
	"fmt"
	"net"
	"time"
)

type MX struct {
	Host string
	IP   []net.IP
	Pref uint16
}

// Obtenir l'enregistrement MX du domaine
// error = nil et map[string]any
// <DOM_RECHERCHE> -> <HOSTS> -> "IP": []net.int , "Pref": int
func (e *EmailValid) GetMX() (json map[string][]MX, err error) {

	defer func() {
		if err := recover(); err != nil {
			json = nil
			err = fmt.Errorf("%s", err)
		}
	}()

	chErr := make(chan error, 1)
	chByte := make(chan map[string][]MX, 1)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Duration(e.timeout.mx)*time.Second)
	defer cancel()

	go func() {

		defer func() {
			if err := recover(); err != nil {
				chErr <- fmt.Errorf("%s", err)
			}
		}()

		mx, err := net.LookupMX(e.domain)
		if err != nil {
			chErr <- err
			return
		}

		listMX := make(map[string][]MX)

		for _, m := range mx {
			ip, _ := net.LookupIP(m.Host)

			listMX[e.domain] = append(listMX[e.domain], MX{
				Host: m.Host,
				IP:   ip,
				Pref: m.Pref,
			})

		}

		chByte <- listMX

	}()

	select {
	case <-ctxTimeout.Done():
		close(chByte)
		return nil, fmt.Errorf("get MX timeout")

	case result := <-chByte:
		close(chErr)
		return result, nil

	case err = <-chErr:
		close(chByte)

	}

	return nil, err
}
