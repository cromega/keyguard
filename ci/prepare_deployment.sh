#!/bin/ash -e

apk update
apk add gettext

buildroot=$PWD
cd code/ci/k8s
envsubst < deployment.yml > temp.yml
mv temp.yml deployment.yml
cd $buildroot

cp code/ci/k8s/* deployment/
echo "$PRIVATE_KEY" > deployment/private_key
