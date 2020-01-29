#!/bin/bash

set -eu -o pipefail

FLY_BIN=${FLY_BIN:-fly}
PIPELINE_NAME="release"

$FLY_BIN validate-pipeline \
	--config "pipelines/release/release.yaml" \
	--yaml-var 'config-approvers=[alphagov]' \
	--yaml-var 'config-approval-count=0' \
	--var "pipeline-name=${PIPELINE_NAME}" \
	--var "branch=validating" \
	--var "github-release-tag-prefix=validating-" \
	"$@"
