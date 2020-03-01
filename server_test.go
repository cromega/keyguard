package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type dummyAuth struct {
	usernameSent, passwordSent string
	called                     int
	ret                        bool
}

func (a *dummyAuth) authenticate(username, password string) (bool, error) {
	a.usernameSent = username
	a.passwordSent = password
	a.called += 1
	return a.ret, nil
}

func TestRootHandlerServesLoaderScriptWithDefaultExpiry(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	server := newServer(
		configuration{LoaderScript: "testdata/loader.sh"},
		nil,
	)
	server.routes()
	server.ServeHTTP(response, request)

	code := response.Code
	if code != 200 {
		t.Error("response code was not 200:", code)
	}

	body := strings.TrimSpace(response.Body.String())
	if body != "awesome loader script 32400" {
		t.Error("wrong response from /:", body)
	}
}

func TestRootHandlerServesLoaderScriptWithExpiry(t *testing.T) {
	request, _ := http.NewRequest("GET", "/3", nil)
	response := httptest.NewRecorder()

	server := newServer(
		configuration{LoaderScript: "testdata/loader.sh"},
		nil,
	)
	server.routes()
	server.ServeHTTP(response, request)

	code := response.Code
	if code != 200 {
		t.Error("response code was not 200:", code)
	}

	body := strings.TrimSpace(response.Body.String())
	if body != "awesome loader script 10800" {
		t.Error("wrong response from /:", body)
	}
}

func TestKeysHandlerRequiresAuthentication(t *testing.T) {
	request, _ := http.NewRequest("GET", "/key", nil)
	response := httptest.NewRecorder()

	server := newServer(
		configuration{PrivateKey: "testdata/id_rsa"},
		&dummyAuth{},
	)
	server.routes()
	server.ServeHTTP(response, request)

	code := response.Code
	if code != 401 {
		t.Error("request should have been rejected with 401:", code)
	}

	header := response.Header().Get("Authenticate")
	if header != "KeyGuard" {
		t.Error("correct authenticate header was not in response:", header)
	}
}

func TestKeysHandlerRequiresValidCredentials(t *testing.T) {
	request, _ := http.NewRequest("GET", "/key", nil)
	response := httptest.NewRecorder()

	request.SetBasicAuth("cromega", "supersecurepassword")

	server := newServer(
		configuration{PrivateKey: "testdata/id_rsa"},
		&dummyAuth{ret: true},
	)
	server.routes()
	server.ServeHTTP(response, request)

	code := response.Code
	if code != 200 {
		t.Error("http status should be 200:", code)
	}
}

func TestKeysHandlerAuthenticatesTheRequest(t *testing.T) {
	request, _ := http.NewRequest("GET", "/key", nil)
	request.SetBasicAuth("keyguard", "supersecurepassword")
	response := httptest.NewRecorder()

	auth := &dummyAuth{ret: true}
	server := newServer(
		configuration{PrivateKey: "testdata/id_rsa"},
		auth,
	)
	server.routes()
	server.ServeHTTP(response, request)

	if auth.called != 1 {
		t.Error("the authenticator was not called")
	}

	if auth.usernameSent != "keyguard" {
		t.Error("sent the wrong username to the authenticator:", auth.usernameSent)
	}

	if auth.passwordSent != "supersecurepassword" {
		t.Error("sent the wrong password to the authenticator:", auth.passwordSent)
	}

	if response.Code != 200 {
		t.Error("wrong respnse code:", response.Code)
	}
}

func TestKeysHandlerRespondsWithKey(t *testing.T) {
	request, _ := http.NewRequest("GET", "/key", nil)
	request.SetBasicAuth("cromega", "supersecurepassword")
	response := httptest.NewRecorder()

	auth := &dummyAuth{ret: true}
	server := newServer(
		configuration{PrivateKey: "testdata/id_rsa"},
		auth,
	)
	server.routes()
	server.ServeHTTP(response, request)

	body := response.Body.String()
	if body != "awesome private key" {
		t.Error("server should have responded with the correct ssh key:", body)
	}
}

func TestPublicKeyHandlerSendsPublicKey(t *testing.T) {
	request, _ := http.NewRequest("GET", "/pubkey", nil)
	response := httptest.NewRecorder()

	server := newServer(
		configuration{PrivateKey: "testdata/real_id_rsa"},
		nil,
	)
	server.routes()
	server.ServeHTTP(response, request)

	if response.Code != 200 {
		t.Error("response should be OK", response.Code)
	}

	body := response.Body.String()
	if body != "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDBgps60FUrGlCFlH48/cSWaDuZYfCTvgRIZt60CiQ1fPj53SQ6xpcDpSCpt1pJt/Q1xZtHPaNZ+HWKAU3tOgspi/AJdrQAPC54CLzdBsMlL/+JxjMxtCf0bbG8dxoRijxIppXVyIuCLabA2oEhepf3U/H+Dvm3XST22f87FsQVrw==\n" {
		t.Error("server should have responded with a public key:", body)
	}
}
