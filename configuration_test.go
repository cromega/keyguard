package main

import (
	"strings"
	"testing"
)

func TestConfigureWorks(t *testing.T) {
	raw := `{"YubiApiKey":"id"}`
	reader := strings.NewReader(raw)

	config, err := configure(reader)

	if err != nil {
		t.Error("config decode failed")
	}

	if config.YubiApiKey != "id" {
		t.Error("config read failed: wrong yubiId: ", config.YubiApiKey)
	}
}

func TestConfigureFailsBecauseOfShitJson(t *testing.T) {
	raw := `{yubiId:id}`
	reader := strings.NewReader(raw)

	_, err := configure(reader)

	if err == nil {
		t.Error("should have failed")
	}
}
