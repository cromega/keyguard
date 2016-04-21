package main

import (
	"fmt"
	"github.com/GeertJohan/yubigo"
)

type yubiAuthenticator struct {
	clientId    string
	apiKey      string
	apiEndpoint string
	preferHttp  bool
}

func (a *yubiAuthenticator) authenticate(username, password string) (bool, error) {
	auth, err := yubigo.NewYubiAuthDebug(a.clientId, a.apiKey, true)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	auth.SetApiServerList(a.apiEndpoint)
	auth.UseHttps(!a.preferHttp)

	_, ok, err := auth.Verify(password)
	if !ok || err != nil {
		fmt.Println(err)
		return false, err
	}

	return true, nil
}
