package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"text/template"
)

type server struct {
	config        configuration
	authenticator authenticator
}

func (s *server) rootHandler(w http.ResponseWriter, r *http.Request) {
	loader, err := readFile(s.config.LoaderScript)
	if err != nil {
		log(err)
		set500Response(w)
		return
	}

	tmpl, err := template.New("loader").Parse(loader)
	if err != nil {
		set500Response(w)
		return
	}

	keyURL := fmt.Sprintf("%s/key", config.PublicURL)
	params := struct {
		URL string
	}{keyURL}
	tmpl.Execute(w, params)
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
		log(err)
		set500Response(w)
		return
	}

	fmt.Fprintf(w, (data))
}

func (s *server) pubKeyHandler(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("ssh-keygen", "-y", "-f", s.config.SSHKey)
	pubkey, err := cmd.Output()

	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			log(string(e.Stderr))
		} else {
			log(e)
		}
		set500Response(w)
		return
	}

	fmt.Fprintf(w, string(pubkey))
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
