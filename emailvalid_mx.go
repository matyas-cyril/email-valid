package emailvalid

import (
	"context"
	libJson "encoding/json"
	"fmt"
	"net"
	"time"
)

// Obtenir l'enregistrement MX du domaine
// error = nil et []byte du json
func (e *EmailValid) GetMX() (json []byte, err error) {

	defer func() {
		if err := recover(); err != nil {
			json = nil
			err = fmt.Errorf("%s", err)
		}
	}()

	chErr := make(chan error, 1)
	chByte := make(chan []byte, 1)

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

		type data struct {
			Name string
			Ip   []net.IP
			Pref uint16
		}

		type host struct {
			Host []data
		}

		h := host{}

		for _, m := range mx {
			ip, _ := net.LookupIP(m.Host)
			d := data{
				Name: m.Host,
				Ip:   ip,
				Pref: m.Pref,
			}
			h.Host = append(h.Host, d)

		}

		j, err := libJson.Marshal(h.Host)
		if err != nil {
			chErr <- err
			return
		}

		chByte <- j

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
