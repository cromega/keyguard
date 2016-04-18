package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestYubiAuthenticatorSendsAuthRequest(t *testing.T) {
	requestBody := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		res, _ := ioutil.ReadAll(r.Body)
		requestBody = string(res)
		fmt.Println(r.URL)
	}))
	defer server.Close()

	auth := yubiAuthenticator{apiEndpoint: server.URL, apiKey: "key"}
	auth.authenticate("password")

	if requestBody != "adsasdads" {
		t.Error("body: ")
	}

}
