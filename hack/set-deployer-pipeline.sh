#!/bin/bash

set -eu

# example of pushing a deployer pipeline for a cluster named CLUSTER_NAME
# using the config at examples/cluster/sandbox.yaml
# but overriding some of the vars
# and setting the git resources to trigger from the current active git branch

PLATFORM_BRANCH=$(git rev-parse --abbrev-ref HEAD)

fly -t cd-gsp sync

fly -t cd-gsp set-pipeline -p "${CLUSTER_NAME}" --config pipelines/deployer/deployer.yaml \
	--load-vars-from pipelines/examples/clusters/sandbox.yaml \
	--var cluster-name=${CLUSTER_NAME} \
	--var cluster-domain=london.${CLUSTER_NAME}.govsvc.uk \
	--var platform-version=${PLATFORM_BRANCH} \
	--var config-version=${PLATFORM_BRANCH} \
	--check-creds

fly -t cd-gsp expose-pipeline -p "${CLUSTER_NAME}"
