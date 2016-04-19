package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type configuration struct {
	YubiApiKey   string
	SSHKey       string
	loaderScript string
}

func configure(r io.Reader) (c configuration, err error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		return
	}

	return
}
