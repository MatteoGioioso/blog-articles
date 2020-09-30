#!/usr/bin/env bash

# Specify an argument to the script to override the environment
# development:   ./deploy.sh
# production:  ./deploy.sh production
# staging:  ./deploy.sh staging
ENV=development
APP_NAME=myapp
PROJECT=${APP_NAME}-${ENV}
BUCKET=${PROJECT}

#Build and testing process
go mod vendor
make build || exit 1

# make the deployment bucket in case it doesn't exist
aws s3 mb s3://"${BUCKET}"

aws cloudformation validate-template \
  --template-body file://template.yaml || exit 1

aws cloudformation package \
  --template-file template.yaml \
  --output-template-file output.yaml \
  --s3-bucket "${BUCKET}"

# the actual deployment step
aws cloudformation deploy \
  --template-file output.yaml \
  --stack-name "${PROJECT}" \
  --capabilities CAPABILITY_IAM \
  --parameter-overrides \
  Environment="$ENV" \
  AppName=${APP_NAME}

