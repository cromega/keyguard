#!/bin/ash -e

apk update
apk add gettext

pushd code/ci/k8s
envsubst < deployment.yml > temp.yml
mv temp.yml deployment.yml
popd

cp code/ci/k8s/* deployment/
