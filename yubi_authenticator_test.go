package main

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestYubiAuthenticatorSendsAuthRequest(t *testing.T) {
	var userid, otp, nonce string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		userid = r.Form["id"][0]
		otp = r.Form["otp"][0]
		nonce = r.Form["nonce"][0]
	}))
	defer server.Close()

	auth := yubiAuthenticator{apiEndpoint: server.URL, apiKey: "key"}
	auth.authenticate("username", "password")

	if userid != "key" {
		t.Error("authenticator sends the wrong key: ", userid)
	}

	if otp != "password" {
		t.Error("authenticator sends the wrong password: ", otp)
	}

	matches, _ := regexp.MatchString("^[0-9a-f]{32}$", nonce)
	if !matches {
		t.Error("authenticator sends the wrong nonce: ", nonce)
	}

}

func TestYubiAuthenticatorReturnsTrueIfResponseIsOk(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	auth := yubiAuthenticator{apiEndpoint: server.URL, apiKey: "key"}
	success, _ := auth.authenticate("username", "password")

	if !success {
		t.Error("auth should have returned true")
	}
}

func TestYubiAuthenticatorReturnsFalseIfResponseIsNotOk(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
	}))
	defer server.Close()

	auth := yubiAuthenticator{apiEndpoint: server.URL, apiKey: "key"}
	success, _ := auth.authenticate("username", "password")

	if success {
		t.Error("auth should have returned false")
	}
}

func TestGetNonceReturnsRandom(t *testing.T) {
	if getNonce() == getNonce() {
		t.Error("nonce is not random")
	}
}
