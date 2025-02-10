package emailvalid

import (
	"context"
	"fmt"
	"net"
	"time"
)

// Obtenir l'enregistrement MX du domaine
// error = nil et map[string]any
// <DOM_RECHERCHE> -> <HOSTS> -> "IP": []net.int , "Pref": int
func (e *EmailValid) GetMX() (json map[string]map[string]any, err error) {

	defer func() {
		if err := recover(); err != nil {
			json = nil
			err = fmt.Errorf("%s", err)
		}
	}()

	chErr := make(chan error, 1)
	chByte := make(chan map[string]map[string]any, 1)

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

		listMX := make(map[string]map[string]any)

		for _, m := range mx {
			ip, _ := net.LookupIP(m.Host)

			// Ajout du premier enregistrement
			if _, exists := listMX[string(e.domain)]; !exists {
				listMX[string(e.domain)] = make(map[string]any)
			}

			listMX[string(e.domain)][m.Host] = map[string]any{
				"IP":   ip,
				"Pref": m.Pref,
			}

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
