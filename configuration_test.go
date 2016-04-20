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

func TestConfigureFallsBackToDefaultValues(t *testing.T) {
	raw := `{"YubiApiKey":"id"}`
	reader := strings.NewReader(raw)

	config, err := configure(reader)

	if err != nil {
		t.Error("config decode failed")
	}

	if config.loaderScript != "loader.sh" {
		t.Error("default settings should have been merged")
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

func TestConfigureValidateFailsIfYubiApiKeyIsMissing(t *testing.T) {
	config = configuration{}

	valid, err := config.validate()
	if valid == true || err == nil {
		t.Error("should have failed")
	}
}

func TestConfigureMergeWithDefaults(t *testing.T) {
	config := configuration{}

	newConfig := mergeWithDefaults(config)

	if newConfig.SSHKey != defaultSSHKey {
		t.Error("SSHKey should have been set to the default value")
	}

	if newConfig.loaderScript != defaultLoaderScript {
		t.Error("loaderScript should have been set to the default value")
	}
}
