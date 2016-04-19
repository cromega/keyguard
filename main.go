package main

import (
	"net/http"
	"os"
)

var config configuration

type server struct {
	config        configuration
	SSHKey        string
	LoaderScript  string
	authenticator authenticator
}

func main() {
	file, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}

	config, err = configure(file)
	if err != nil {
		panic(err)
	}

	server := server{config: config}
	http.HandleFunc("/", server.rootHandler)
	http.HandleFunc("/keys", server.keysHandler)
}
