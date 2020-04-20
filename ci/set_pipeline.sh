#!/bin/bash -e

lpass status || { echo "log in to lastpass first."; exit 1; }

fly -t ci set-pipeline -p keyguard -c ci/pipeline.yml \
  -v "deploy-key=$(lpass show keyguard-deploy-key --field='Private Key')" \
  -v "dockerhub-password=$(lpass show --password docker.com)"
