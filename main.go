package main

import (
	logger "log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

type configuration struct {
	PrivateKey   string `envconfig:"PRIVATE_KEY" default:"id_rsa"`
	LoaderScript string `envconfig:"LOADER_SCRIPT" default:"loader.sh"`
	PublicURL    string `envconfig:"PUBLIC_URL" required:"true"`
	AuthModule   string `envconfig:"AUTH_MODULE" default:"yubikey"`
	Port         string `envconfig:"PORT" default:"8000"`
}

func main() {
	server, err := initialize()
	if err != nil {
		logger.Fatal(err)
	}

	http.ListenAndServe(":"+server.config.Port, server.router)
}

func initialize() (server, error) {
	var c configuration
	err := envconfig.Process("KG", &c)
	if err != nil {
		envconfig.Usage("KG", &c)
		logger.Fatal(err)
	}

	var auth authenticator
	switch c.AuthModule {
	case "yubikey":
		auth, err = NewYubiAuthenticator()
		if err != nil {
			logger.Fatal(err)
		}
	default:
		logger.Fatal("no such authenticator: " + c.AuthModule)
	}

	server := newServer(c, auth)
	server.routes()

	return server, nil
}

func log(message interface{}) {
	logger.Println(message)
}
