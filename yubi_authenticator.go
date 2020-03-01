package main

import (
	"github.com/GeertJohan/yubigo"
	"github.com/kelseyhightower/envconfig"
)

type yubiConfig struct {
	ClientId string `envconfig:"CLIENT_ID" required:"true"`
	ApiKey   string `envconfig:"API_KEY" required:"true"`
	ApiHost  string `envconfig:"API_HOST" default:"api.yubico.com/wsapi/2.0/verify"`
	UseHttps bool   `envconfig:"USE_HTTPS" default:"true"`
}

type yubiAuth struct {
	authenticator *yubigo.YubiAuth
}

func NewYubiAuthenticator() (authenticator, error) {
	var c yubiConfig
	err := envconfig.Process("KG_YUBI", &c)
	if err != nil {
		envconfig.Usage("KG_YUBI", &c)
		return nil, err
	}

	yubi, err := yubigo.NewYubiAuth(c.ClientId, c.ApiKey)
	if err != nil {
		return nil, err
	}

	yubi.SetApiServerList(c.ApiHost)
	yubi.UseHttps(c.UseHttps)

	auth := yubiAuth{authenticator: yubi}

	return &auth, nil
}

func (a *yubiAuth) authenticate(_, password string) (ok bool, err error) {
	_, ok, err = a.authenticator.Verify(password)
	return
}
