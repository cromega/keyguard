package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
)

type configuration struct {
	YubiApiKey   string
	SSHKey       string
	loaderScript string
}

const (
	defaultSSHKey       = "id_rsa"
	defaultLoaderScript = "loader.sh"
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

func (c *configuration) validate() (bool, error) {
	if c.YubiApiKey == "" {
		return false, errors.New("YubiApiKey is missing from configuration")
	}

	return true, nil
}

func mergeWithDefaults(c configuration) configuration {
	if c.SSHKey == "" {
		c.SSHKey = defaultSSHKey
	}

	if c.loaderScript == "" {
		c.loaderScript = defaultLoaderScript
	}

	return c
}
