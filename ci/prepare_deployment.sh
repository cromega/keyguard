#!/bin/ash -e

apk update
apk add gettext

buildroot=$PWD
cd code/ci/k8s

for file in *.yml; do
  envsubst < $file > temp
  mv temp $file
done

cd $buildroot
cp code/ci/k8s/* deployment/
