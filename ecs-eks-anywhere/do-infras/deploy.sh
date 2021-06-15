#!/usr/bin/env sh

aws ssm create-activation --iam-role ECSAnywhereRole --registration-limit 1 | tee ssm-activation.json
terraform apply -auto-approve -var-file secret.tfvars