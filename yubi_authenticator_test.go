package main

import (
	"os"
	"testing"
)

func TestNewYubiAuthenticator(t *testing.T) {
	os.Setenv("KG_YUBI_CLIENT_ID", "123")
	os.Setenv("KG_YUBI_API_KEY", "a2V5")

	_, err := NewYubiAuthenticator()

	if err != nil {
		t.Error("creating authenticator failed: ", err)
	}
}
