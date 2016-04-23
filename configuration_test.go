package main

import (
	"strings"
	"testing"
)

func TestConfigureWorks(t *testing.T) {
	raw := `{"loaderScript":"loader.sh"}`
	reader := strings.NewReader(raw)

	config, err := configure(reader)

	if err != nil {
		t.Error("config decode failed")
	}

	if config.loaderScript != "loader.sh" {
		t.Error("config read failed: wrong yubiId: ", config.loaderScript)
	}
}

func TestConfigureFallsBackToDefaultValues(t *testing.T) {
	raw := `{}`
	reader := strings.NewReader(raw)

	config, err := configure(reader)

	if err != nil {
		t.Error("config decode failed")
	}

	if config.loaderScript != defaultLoaderScript {
		t.Error("default loader script setting should have been merged: ", config.loaderScript)
	}

	if config.SSHKey != defaultSSHKey {
		t.Error("default sshkey setting should have been merged: ", config.SSHKey)
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
