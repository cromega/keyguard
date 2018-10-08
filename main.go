package main

import (
	"flag"
	logger "log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

var config configuration

func main() {
	var configPath = flag.String("configPath", "config.json", "path to the config file")
	flag.Parse()

	file, err := os.Open(*configPath)
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
	router := httprouter.New()
	router.GET("/", server.rootHandler)
	router.GET("/key", server.keyHandler)
	router.GET("/pubkey", server.pubKeyHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3459"
	}

	http.ListenAndServe(":"+port, router)
}

func log(message interface{}) {
	logger.Println(message)
}
