package main

type authenticator interface {
	authenticate(username, password string) (ok bool, err error)
}
