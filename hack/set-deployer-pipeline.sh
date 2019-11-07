#!/bin/bash

set -eu -o pipefail

: "${CLUSTER_CONFIG:?}"

FLY_BIN=${FLY_BIN:-fly}
CLUSTER_NAME=$(yq -r '.["cluster-name"]' < ${CLUSTER_CONFIG})
PIPELINE_NAME=$(yq -r '.["concourse-pipeline-name"]' < ${CLUSTER_CONFIG})

echo "generating approvers for ${CLUSTER_NAME}..."


$FLY_BIN -t cd-gsp sync

$FLY_BIN -t cd-gsp set-pipeline -p "${PIPELINE_NAME}" \
	--config "pipelines/deployer/deployer.yaml" \
	--load-vars-from "pipelines/deployer/deployer.defaults.yaml" \
	--load-vars-from "${CLUSTER_CONFIG}" \
	--yaml-var 'config-approvers=[alphagov]' \
	--yaml-var 'config-approval-count=0' \
	--yaml-var 'trusted-developer-keys=[]' \
	--check-creds "$@"

$FLY_BIN -t cd-gsp expose-pipeline -p "${PIPELINE_NAME}"

