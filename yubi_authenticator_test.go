package main

import (
	"testing"
)

func TestNewAuthenticator(t *testing.T) {
	config := map[string]interface{}{
		"clientId": "id",
		"apiKey":   "a2V5",
	}

	_, err := NewAuthenticator(config)

	if err != nil {
		t.Error("creating authenticator failed: ", err)
	}
}
