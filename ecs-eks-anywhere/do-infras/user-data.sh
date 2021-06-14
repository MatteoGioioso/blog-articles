#!/usr/bin/env sh

# Register instance
ACTIVATION_ID=$(jq '.ActivationId' < ssm-activation.json)
ACTIVATION_CODE=$(jq '.ActivationCode' < ssm-activation.json)

export ACTIVATION_ID=${ACTIVATION_ID}
export ACTIVATION_CODE=${ACTIVATION_CODE}
export REGION=us-east-1
export CLUSTER_NAME=ecs-anywhere-test-cluster

curl --proto \"https\" -o \"/tmp/ecs-anywhere-install.sh\" \"https://amazon-ecs-agent.s3.amazonaws.com/ecs-anywhere-install-latest.sh\"
sudo bash /tmp/ecs-anywhere-install.sh --region "$REGION" --cluster "$CLUSTER_NAME" --activation-id "$ACTIVATION_ID" --activation-code "$ACTIVATION_CODE"
