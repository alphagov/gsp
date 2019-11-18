#!/bin/bash

set -eu -o pipefail

: "${USER_CONFIGS:?}"

FLY_BIN=${FLY_BIN:-fly}
PIPELINE_NAME="release"

CURRENT_BRANCH="$(git rev-parse --abbrev-ref HEAD)"
if [ "${CURRENT_BRANCH}" != "master" ]
then
	echo "${CURRENT_BRANCH} is not master!"
	exit 1
fi

echo "generating initial list of trusted developers for releases..."

approvers="/tmp/gsp-release-approvers.yaml"
echo -n "config-approvers: " > "${approvers}"
yq . ${USER_CONFIGS}/*.yaml \
	| jq -c -s "[.[].github] | unique | sort" \
	>> "${approvers}"

$FLY_BIN -t cd-gsp sync

$FLY_BIN -t cd-gsp set-pipeline -p "${PIPELINE_NAME}" \
	--config "pipelines/release/release.yaml" \
	--load-vars-from "${approvers}" \
	--var "pipeline-name=${PIPELINE_NAME}" \
	--var "branch=${CURRENT_BRANCH}" \
	--var "github-release-tag-prefix=gsp-" \
	--check-creds "$@"

$FLY_BIN -t cd-gsp expose-pipeline -p "${PIPELINE_NAME}"
