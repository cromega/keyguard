package main

import (
	"net/http"
	"net/http/httptest"
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

func TestRootHandlerServesLoaderScript(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "", nil)

	server := server{config: configuration{loaderScript: "loader script"}}
	server.rootHandler(response, request)

	code := response.Code
	if code != 200 {
		t.Error("response code was not 200: ", code)
	}

	body := response.Body.String()
	if response.Body.String() != "loader script" {
		t.Error("wrong response from / endpoint: ", body)
	}
}

func TestKeysHandlerRequiresAuthentication(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "", nil)

	server := server{authenticator: &dummyAuth{}}
	server.keysHandler(response, request)

	code := response.Code
	if code != 401 {
		t.Error("request should have been rejected with 401: ", code)
	}

	header := response.Header().Get("Authenticate")
	if header != "KeyGuard" {
		t.Error("correct authenticate header was not in response: ", header)
	}
}

func TestKeysHandlerRequiresValidCredentials(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "", nil)

	request.SetBasicAuth("cromega", "supersecurepassword")

	server := server{authenticator: &dummyAuth{ret: true}}
	server.keysHandler(response, request)

	code := response.Code
	if code != 200 {
		t.Error("request should have been accepted")
	}
}

func TestKeysHandlerAuthenticatesTheRequest(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "", nil)

	request.SetBasicAuth("cromega", "supersecurepassword")

	auth := &dummyAuth{ret: true}
	server := server{authenticator: auth}
	server.keysHandler(response, request)

	if auth.called != 1 {
		t.Error("the authenticator was not called")
	}

	if auth.usernameSent != "cromega" {
		t.Error("sent the wrong username to the authenticator: ", auth.usernameSent)
	}

	if auth.passwordSent != "supersecurepassword" {
		t.Error("sent the wrong password to the authenticator: ", auth.passwordSent)
	}
}

func TestKeysHandlerRespondsWithKey(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "", nil)

	request.SetBasicAuth("cromega", "supersecurepassword")

	server := server{config: configuration{SSHKey: "ssh key"}, authenticator: &dummyAuth{ret: true}}
	server.keysHandler(response, request)

	body := response.Body.String()
	if body != "ssh key" {
		t.Error("server should have responded with the correct ssh key: ", body)
	}
}
