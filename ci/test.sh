#!/bin/ash -e

export CGO_ENABLED=0

apk update
apk add openssh-keygen

cd code
chmod 600 testdata/real_id_rsa
go test -v
