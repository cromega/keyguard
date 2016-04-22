package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func (s *server) rootHandler(w http.ResponseWriter, r *http.Request) {
	data, err := readFile(s.config.loaderScript)
	if err != nil {
		set500Response(w)
		return
	}

	fmt.Fprintf(w, (data))
}

func (s *server) keyHandler(w http.ResponseWriter, r *http.Request) {
	username, password, _ := r.BasicAuth()

	authenticated, err := s.authenticator.authenticate(username, password)
	if err != nil {
		w.WriteHeader(503)
		return
	}

	if !authenticated {
		set401Response(w)
		return
	}

	data, err := readFile(s.config.SSHKey)
	if err != nil {
		set500Response(w)
		return
	}

	fmt.Fprintf(w, (data))
}

func set401Response(w http.ResponseWriter) {
	w.Header().Add("Authenticate", "KeyGuard")
	w.WriteHeader(401)
}

func set500Response(w http.ResponseWriter) {
	w.WriteHeader(500)
}

func readFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	data, _ := ioutil.ReadAll(f)
	return string(data), nil
}
