package main

import (
	"fmt"
	"net/http"
)

func (s *server) rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.loaderScript)
}

func (s *server) keysHandler(w http.ResponseWriter, r *http.Request) {
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

	fmt.Fprintf(w, s.config.SSHKey)
}

func set401Response(w http.ResponseWriter) {
	w.Header().Add("Authenticate", "KeyGuard")
	w.WriteHeader(401)
}
