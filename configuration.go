package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type configuration struct {
	SSHKey       string
	loaderScript string
	auth         map[string]interface{}
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

func mergeWithDefaults(c configuration) configuration {
	if c.SSHKey == "" {
		c.SSHKey = defaultSSHKey
	}

	if c.loaderScript == "" {
		c.loaderScript = defaultLoaderScript
	}

	return c
}
