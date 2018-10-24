package main

import (
	"errors"
	"github.com/GeertJohan/yubigo"
)

type yubiAuth struct {
	authenticator *yubigo.YubiAuth
}

const (
	defaultApiEndpoint = "api.yubico.com/wsapi/2.0/verify"
)

func NewAuthenticator(config map[string]interface{}) (authenticator, error) {
	raw := config["clientId"]
	clientId, ok := raw.(string)
	if !ok || clientId == "" {
		return nil, errors.New("missing clientId from auth config")
	}

	raw = config["apiKey"]
	apiKey, ok := raw.(string)
	if !ok || apiKey == "" {
		return nil, errors.New("missing clientId from auth config")
	}

	yubi, err := yubigo.NewYubiAuth(clientId, apiKey)
	if err != nil {
		return nil, err
	}

	auth := yubiAuth{authenticator: yubi}

	raw = config["apiEndpoint"]
	if apiEndpoint, _ := raw.(string); apiEndpoint != "" {
		auth.authenticator.SetApiServerList(apiEndpoint)
	}

	raw = config["preferHttp"]
	preferHttp, _ := raw.(bool)
	auth.authenticator.UseHttps(!preferHttp)

	return &auth, nil
}

func (a *yubiAuth) authenticate(_, password string) (ok bool, err error) {
	_, ok, err = a.authenticator.Verify(password)
	return
}
