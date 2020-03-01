#!/bin/bash -e

fly -t devbox set-pipeline -p keyguard -c ci/pipeline.yml \
  -v "dockerhub-password=$(lpass show --password docker.com)" \
  -v "deploy-key=$(lpass show keyguard-deploy-key --field='Private Key')" \
  -v "ci-access-key-id=$(lpass show ci-user-access-key --field=access-key-id)" \
  -v "ci-secret-access-key=$(lpass show ci-user-access-key --field=secret-access-key)" \
  -v "yubico-client-id=$(lpass show yubico --field=client-id)" \
  -v "yubico-api-key=$(lpass show yubico --field=api-key)" \
  -v "private-key=$(lpass show devkey --field='Private Key' | base64 -w0)"
