#!/bin/bash

set -eu -o pipefail

: "${CLUSTER_CONFIG:?}"

CLUSTER_NAME=$(yq -r '.["cluster-name"]' < ${CLUSTER_CONFIG})
PIPELINE_NAME=$(yq -r '.["concourse-pipeline-name"]' < ${CLUSTER_CONFIG})
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

echo ""

fly -t cd-gsp sync

fly -t cd-gsp set-pipeline -p "${PIPELINE_NAME}" \
	--config "pipelines/deployer/deployer.yaml" \
	--load-vars-from "pipelines/deployer/deployer.defaults.yaml" \
	--load-vars-from "${CLUSTER_CONFIG}" \
	--var "platform-version=${CURRENT_BRANCH}" \
	--var "config-version=${CURRENT_BRANCH}" \
	--var "users-version=${CURRENT_BRANCH}" \
	--check-creds

fly -t cd-gsp expose-pipeline -p "${PIPELINE_NAME}"

