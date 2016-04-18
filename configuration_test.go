package main

import (
	"strings"
	"testing"
)

func TestConfigureWorks(t *testing.T) {
	raw := `{"YubiId":"id"}`
	reader := strings.NewReader(raw)

	config, err := configure(reader)

	if err != nil {
		t.Error("config decode failed")
	}

	if config.YubiId != "id" {
		t.Error("config read failed: wrong yubiId: ", config.YubiId)
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
