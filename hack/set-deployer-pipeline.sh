#!/bin/bash

set -eu -o pipefail

: "${CLUSTER_CONFIG:?}"
: "${USER_CONFIGS:?}"

CLUSTER_NAME=$(yq -r '.["cluster-name"]' < ${CLUSTER_CONFIG})
PIPELINE_NAME=$(yq -r '.["concourse-pipeline-name"]' < ${CLUSTER_CONFIG})
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

if [[ "${CURRENT_BRANCH}" != "master" ]]; then
	echo "WARNING: current branch is not master!!!"
	read -p "Are you sure you want to set pipeline with platform-version pinned to ${CURRENT_BRANCH}? " -n 1 -r
	echo
	if [[ ! $REPLY =~ ^[Yy]$ ]]
	then
		[[ "$0" = "$BASH_SOURCE" ]] && exit 1 || return 1
	fi
fi

approvers="/tmp/deployer-${CLUSTER_NAME}-approvers.yaml"
echo -n "github-approvers: " > "${approvers}"
cat ${USER_CONFIGS}/*.yaml \
	| yq . \
	| jq -c -s "[.[] | select(.roles[] | select((. == \"${CLUSTER_NAME}-sre\" ) or (. == \"${CLUSTER_NAME}-admin\"))) | .github] | unique | sort" \
	>> "${approvers}"

trusted="/tmp/deployer-${CLUSTER_NAME}-keys.yaml"
echo -n "trusted-developer-keys: " > "${trusted}"
cat ${USER_CONFIGS}/*.yaml \
	| yq . \
	| jq -c -s '[ .[].pub ] | sort' \
	>> "${trusted}"

fly -t cd-gsp sync

fly -t cd-gsp set-pipeline -p "${PIPELINE_NAME}" \
	--config "pipelines/deployer/deployer.yaml" \
	--load-vars-from "pipelines/deployer/deployer.defaults.yaml" \
	--load-vars-from "${CLUSTER_CONFIG}" \
	--var "platform-version=${CURRENT_BRANCH}" \
	--load-vars-from "${approvers}" \
	--load-vars-from "${trusted}" \
	--check-creds "$@"

fly -t cd-gsp expose-pipeline -p "${PIPELINE_NAME}"

