#!/bin/bash

set -eu -o pipefail

FLY_BIN=${FLY_BIN:-fly}
CLUSTER_NAME="validation"
PIPELINE_NAME="validation-deployer"

$FLY_BIN validate-pipeline \
	--config "pipelines/deployer/deployer.yaml" \
	--load-vars-from "pipelines/deployer/deployer.defaults.yaml" \
	--yaml-var 'config-approvers=[alphagov]' \
	--yaml-var 'config-approval-count=0' \
	"$@"


