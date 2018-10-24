package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"text/template"

	"github.com/AndrewBurian/powermux"
)

type server struct {
	config        configuration
	authenticator authenticator
	router        *powermux.ServeMux
}

func newServer(c configuration, a authenticator) server {
	return server{
		config:        c,
		authenticator: a,
		router:        powermux.NewServeMux(),
	}
}

const (
	DefaultExpiry = 9
)

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) rootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loader, err := readFile(s.config.LoaderScript)
		if err != nil {
			log(err)
			setErrorResponse(500, w)
			return
		}

		tmpl, err := template.New("loader").Parse(loader)
		if err != nil {
			log(err)
			setErrorResponse(500, w)
			return
		}

		expiry := DefaultExpiry
		if expiryParam := powermux.PathParam(r, "expiry"); expiryParam != "" {
			expiry, _ = strconv.Atoi(expiryParam)
		}

		keyURL := fmt.Sprintf("%s/key", config.PublicURL)
		expirySeconds := expiry * 3600

		templateParams := struct {
			URL    string
			Expiry int
		}{keyURL, expirySeconds}
		tmpl.Execute(w, templateParams)
	}
}

func (s *server) keyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, _ := r.BasicAuth()

		authenticated, err := s.authenticator.authenticate(username, password)
		if err != nil {
			log(err)
			setErrorResponse(503, w)
			return
		}

		if !authenticated {
			set401Response(w)
			return
		}

		data, err := readFile(s.config.SSHKey)
		if err != nil {
			log(err)
			setErrorResponse(500, w)
			return
		}

		fmt.Fprintf(w, (data))
	}
}

func (s *server) pubKeyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cmd := exec.Command("ssh-keygen", "-y", "-f", s.config.SSHKey)
		pubkey, err := cmd.Output()

		if err != nil {
			if e, ok := err.(*exec.ExitError); ok {
				log(string(e.Stderr))
			} else {
				log(e)
			}
			setErrorResponse(500, w)
			return
		}

		fmt.Fprintf(w, string(pubkey))
	}
}

func set401Response(w http.ResponseWriter) {
	w.Header().Add("Authenticate", "KeyGuard")
	w.WriteHeader(401)
}

func setErrorResponse(code int, w http.ResponseWriter) {
	w.WriteHeader(code)
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
