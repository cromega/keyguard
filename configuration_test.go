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

	if config.LoaderScript != "loader.sh" {
		t.Error("config read failed: wrong yubiId: ", config.LoaderScript)
	}
}

func TestConfigureFallsBackToDefaultValues(t *testing.T) {
	raw := `{}`
	reader := strings.NewReader(raw)

	config, err := configure(reader)

	if err != nil {
		t.Error("config decode failed")
	}

	if config.LoaderScript != defaultLoaderScript {
		t.Error("default loader script setting should have been merged: ", config.LoaderScript)
	}

	if config.SSHKey != defaultSSHKey {
		t.Error("default sshkey setting should have been merged: ", config.SSHKey)
	}

	if config.SSHPubKey != defaultSSHPubKey {
		t.Error("default ssh public key setting should have been merged: ", config.SSHPubKey)
	}
}

func TestConfigureMergeWithDefaults(t *testing.T) {
	config := configuration{}

	newConfig := mergeWithDefaults(config)

	if newConfig.SSHKey != defaultSSHKey {
		t.Error("SSHKey should have been set to the default value")
	}

	if newConfig.SSHPubKey != defaultSSHPubKey {
		t.Error("SSHPubKey should have been set to the default value")
	}

	if newConfig.LoaderScript != defaultLoaderScript {
		t.Error("loaderScript should have been set to the default value")
	}
}
