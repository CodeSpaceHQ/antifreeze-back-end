#!/bin/bash

# This file is for Travis CI

docker login -u "$DOCKER_USER" -p "$DOCKER_PASS"
docker build -f deploy/docker/prod.Dockerfile -t nilsgs/antifreeze:latest -t nilsgs/antifreeze:2.9 .
docker push nilsgs/antifreeze:2.9
docker push nilsgs/antifreeze:latest

echo $SERVICE_ACCOUNT | base64 --decode > $GOOGLE_APPLICATION_CREDENTIALS
echo $TF_SECRETS | base64 --decode > deploy/terraform/secret.tfvars

cd deploy/terraform
terraform init
terraform apply -auto-approve --var-file="secret.tfvars"

cd $TRAVIS_BUILD_DIR
rm $GOOGLE_APPLICATION_CREDENTIALS
rm deploy/terraform/secret.tfvars
