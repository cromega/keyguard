package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
)

type authenticator interface {
	authenticate(username, password string) (ok bool, err error)
}

type yubiAuthenticator struct {
	apiEndpoint string
	apiKey      string
}

func (a *yubiAuthenticator) authenticate(username, password string) (ok bool, err error) {
	url := fmt.Sprintf("%s?id=%s&otp=%s&nonce=%s", a.apiEndpoint, a.apiKey, password, getNonce())

	response, err := http.Get(url)
	if err != nil {
		return false, err
	}

	if response.StatusCode != 200 {
		return false, nil
	}

	ok = true
	return
}

func getNonce() (nonce string) {
	bytes := make([]byte, 16)
	rand.Read(bytes)

	for _, b := range bytes {
		nonce = nonce + fmt.Sprintf("%02x", b)
	}

	return
}
