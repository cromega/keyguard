package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
)

type authenticator interface {
	authenticate(password string) (ok bool, err error)
}

type yubiAuthenticator struct {
	apiEndpoint string
	apiKey      string
}

func (a *yubiAuthenticator) authenticate(password string) (ok bool, err error) {
	url := fmt.Sprintf("%s?id=%s&otp=%s&nonce=%s", a.apiEndpoint, a.apiKey, password, getNonce())

	_, err = http.Get(url)
	if err != nil {
		return false, err
	}

	ok = true
	return
}

func getNonce() (nonce string) {
	bytes := make([]byte, 16)
	rand.Read(bytes)

	for b := range bytes {
		nonce = nonce + fmt.Sprintf("%02x", b)
	}

	return
}
