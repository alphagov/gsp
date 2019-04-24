#!/bin/bash

set -eu

PLATFORM_BRANCH=$(git rev-parse --abbrev-ref HEAD)
PLATFORM_REPO="https://github.com/alphagov/gsp-terraform-ignition.git"

SANDBOX_ACCOUNT_ID="011571571136"

fly4 -t cd-gsp set-pipeline -p farms --config pipelines/deployer/deployer.yaml \
	--var account-id=$SANDBOX_ACCOUNT_ID \
	--var account-name=sandbox \
	--var account-role-arn=arn:aws:iam::$SANDBOX_ACCOUNT_ID:role/deployer \
	--var cluster-name=${CLUSTER_NAME} \
	--yaml-var trusted-developer-keys="[]" \
	--var splunk-enabled=0 \
	--var splunk-hec-token="NOTATOKEN" \
	--var splunk-hec-url=NOTAURL \
	--var github-client-secret=NOTASECRET \
	--var github-client-id=NOTID \
	--var eks-version=1.12 \
	--var platform-repository=${PLATFORM_REPO} \
	--var platform-version=${PLATFORM_BRANCH} \
	--var config-repository=${PLATFORM_REPO} \
	--var config-version=${PLATFORM_BRANCH} \
	--var config-path=pipelines/examples \
	--check-creds

