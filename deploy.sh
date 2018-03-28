#!/bin/bash

# This file is for Travis CI

docker login -u "$DOCKER_USER" -p "$DOCKER_PASS"
docker build -f .\deploy\docker\prod.Dockerfile -t nilsgs/antifreeze:latest -t nilsgs/antifreeze:3.0 .
docker push nilsgs/antifreeze:3.0
docker push nilsgs/antifreeze:latest

echo $SERVICE_ACCOUNT | base64 --decode > key.json
echo $TF_SECRETS | base64 --decode > deploy/terraform/secret.tfvars

cd deploy/terraform
terraform init
terraform apply -auto-approve --var-file="secret.tfvars"

cd $TRAVIS_BUILD_DIR
rm key.json
rm deploy/terraform/secret.tfvars
