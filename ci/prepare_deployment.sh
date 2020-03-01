#!/bin/ash -e

apk update
apk add gettext

buildroot=$PWD
cd code/ci/k8s
envsubst < deployment.yml > temp.yml
mv temp.yml deployment.yml

export PRIVATE_KEY_BASE64="$(echo "$PRIVATE_KEY" | base64)"
envsubst < secret.yml > temp.yml
mv temp.yml secret.yml
cd $buildroot

cp code/ci/k8s/* deployment/
