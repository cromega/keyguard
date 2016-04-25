package main

import (
	"net/http"
	"os"
)

var config configuration

func main() {
	file, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	config, err = configure(file)
	if err != nil {
		panic(err)
	}

	auth, err := NewAuthenticator(config.Auth)

	if err != nil {
		panic(err)
	}

	server := server{config: config, authenticator: auth}
	http.HandleFunc("/", server.rootHandler)
	http.HandleFunc("/key", server.keyHandler)

	http.ListenAndServe(":3456", nil)
}
