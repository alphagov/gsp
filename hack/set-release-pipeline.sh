#!/bin/bash

set -eu -o pipefail

: "${USER_CONFIGS:?}"

PIPELINE_NAME="release"

echo "generating initial list of trusted developers for releases..."

approvers="/tmp/gsp-release-approvers.yaml"
echo -n "github-approvers: " > "${approvers}"
yq . ${USER_CONFIGS}/*.yaml \
	| jq -c -s "[.[].github] | unique | sort" \
	>> "${approvers}"

trusted="/tmp/gsp-release-keys.yaml"
echo -n "trusted-developer-keys: " > "${trusted}"
yq . ${USER_CONFIGS}/*.yaml \
	| jq -c -s '[ .[].pub ] | sort' \
	>> "${trusted}"

fly -t cd-gsp sync

fly -t cd-gsp set-pipeline -p "${PIPELINE_NAME}" \
	--config "pipelines/release/release.yaml" \
	--load-vars-from "${approvers}" \
	--load-vars-from "${trusted}" \
	--var "pipeline-name=${PIPELINE_NAME}" \
	--check-creds "$@"

fly -t cd-gsp expose-pipeline -p "${PIPELINE_NAME}"

