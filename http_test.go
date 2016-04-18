package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type dummyAuthenticator struct {
}

func TestRootHandlerServesLoaderScript(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "", nil)

	server := server{LoaderScript: "loader script"}
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

	server := server{}
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

	request.SetBasicAuth("crome", "supersecurepassword")

	server := server{}
	server.keysHandler(response, request)

	code := response.Code
	if code != 200 {
		t.Error("request should have been accepted")
	}
}

func TestKeysHandlerRespondsWithKey(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "", nil)

	request.SetBasicAuth("crome", "supersecurepassword")

	server := server{SSHKey: "ssh key"}
	server.keysHandler(response, request)

	body := response.Body.String()
	if body != "ssh key" {
		t.Error("server should have responded with the correct ssh key: ", body)
	}
}
