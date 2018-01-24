package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type configuration struct {
	SSHKey       string
	SSHPubKey    string
	LoaderScript string
	PublicURL    string
	Auth         map[string]interface{}
}

const (
	defaultSSHKey       = "id_rsa"
	defaultLoaderScript = "loader.sh"
	defaultSSHPubKey    = "id_rsa.pub"
)

func configure(r io.Reader) (c configuration, err error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		return
	}

	return mergeWithDefaults(c), nil
}

func mergeWithDefaults(c configuration) configuration {
	if c.SSHKey == "" {
		c.SSHKey = defaultSSHKey
	}

	if c.SSHPubKey == "" {
		c.SSHPubKey = defaultSSHPubKey
	}

	if c.LoaderScript == "" {
		c.LoaderScript = defaultLoaderScript
	}

	return c
}
