package emailvalid

import (
	"fmt"
)

// Déclarer un objet de Email valide
func NewEmail(data string) (email *EmailValid, err error) {

	defer func() {
		if errRecover := recover(); errRecover != nil {
			email = nil
			err = fmt.Errorf("%s", err)
		}
	}()

	email, err = extractEmail(data)
	if err != nil {
		return nil, err
	}

	return email, nil
}

func (e *EmailValid) GetAllEmailData() (data map[string]string, err error) {

	defer func() {
		if errRecover := recover(); errRecover != nil {
			data = nil
			err = fmt.Errorf("%s", err)
		}
	}()

	e.mu.Lock()
	defer e.mu.Unlock()

	return map[string]string{
		"SrcEmail": e.srcEmail,
		"Email":    e.email,
		"Name":     e.name,
		"Local":    e.local,
		"Domain":   e.domain,
		"TLD":      e.domArrayInv[0],
	}, nil

}

func (e *EmailValid) GetEmail() (string, error) {

	data, err := e.GetAllEmailData()
	if err != nil {
		return "", err
	}
	return data["Email"], nil

}

func (e *EmailValid) GetDomain() (string, error) {

	data, err := e.GetAllEmailData()
	if err != nil {
		return "", err
	}
	return data["Domain"], nil
}

func (e *EmailValid) GetLocal() (string, error) {

	data, err := e.GetAllEmailData()
	if err != nil {
		return "", err
	}
	return data["Local"], nil
}

func (e *EmailValid) GetTLD() (string, error) {

	data, err := e.GetAllEmailData()
	if err != nil {
		return "", err
	}

	if len(data["TLD"]) == 0 {
		return "", fmt.Errorf("no TLD available")
	}

	return data["TLD"], nil
}

func (e *EmailValid) GetTimeoutMX() (timeout uint, err error) {

	defer func() {
		if errRecover := recover(); errRecover != nil {
			timeout = 0
			err = fmt.Errorf("%s", err)
		}
	}()

	e.mu.Lock()
	defer e.mu.Unlock()

	return e.timeout.mx, nil
}

func (e *EmailValid) SetTimeoutMX(timeout uint) (err error) {

	defer func() {
		if errRecover := recover(); errRecover != nil {
			err = fmt.Errorf("%s", err)
		}
	}()

	e.mu.Lock()
	defer e.mu.Unlock()

	e.timeout.mx = timeout
	return nil

}

func (e *EmailValid) GetTimeoutSPF() (timeout uint, err error) {

	defer func() {
		if errRecover := recover(); errRecover != nil {
			timeout = 0
			err = fmt.Errorf("%s", err)
		}
	}()

	e.mu.Lock()
	defer e.mu.Unlock()

	return e.timeout.spf, nil
}

func (e *EmailValid) SetTimeoutSPF(timeout uint) (err error) {

	defer func() {
		if errRecover := recover(); errRecover != nil {
			err = fmt.Errorf("%s", err)
		}
	}()

	e.mu.Lock()
	defer e.mu.Unlock()

	e.timeout.spf = timeout
	return nil

}

func (e *EmailValid) GetTimeoutDMARC() (timeout uint, err error) {

	defer func() {
		if errRecover := recover(); errRecover != nil {
			timeout = 0
			err = fmt.Errorf("%s", err)
		}
	}()

	e.mu.Lock()
	defer e.mu.Unlock()

	return e.timeout.dmarc, nil
}

func (e *EmailValid) SetTimeoutDMARC(timeout uint) (err error) {

	defer func() {
		if errRecover := recover(); errRecover != nil {
			err = fmt.Errorf("%s", err)
		}
	}()

	e.mu.Lock()
	defer e.mu.Unlock()

	e.timeout.dmarc = timeout
	return nil

}
