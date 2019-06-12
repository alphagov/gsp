#!/usr/bin/env bash

set -eu
set -o pipefail

: "${CODECOMMIT_REPO_URL:?Required variable!}"
: "${SOURCE_REPO_URL:?Required variable!}"

temp_dir="$(mktemp -d)"
pushd "$temp_dir"
cleanup() {
  popd
  rm -rf "$temp_dir"
}
trap cleanup EXIT

if [ -n "${CODECOMMIT_INIT_ROLE_ARN}" ]
then
    echo "Assuming role ${CODECOMMIT_INIT_ROLE_ARN}..."
	temp_role="$(aws sts assume-role --role-arn "${CODECOMMIT_INIT_ROLE_ARN}" --role-session-name "terraform-init-codecommit" | jq .Credentials)"
	export AWS_ACCESS_KEY_ID="$(echo $temp_role | jq -r .AccessKeyId)"
	export AWS_SECRET_ACCESS_KEY="$(echo $temp_role | jq -r .SecretAccessKey)"
	export AWS_SESSION_TOKEN="$(echo $temp_role | jq -r .SessionToken)"

	# Looks like aws-vault sets this env variable and it screws with the client somehow
	unset AWS_SECURITY_TOKEN
fi

git init
git config --local credential.helper '!aws codecommit credential-helper $@'
git config --local credential.UseHttpPath true
git config --local user.name "Friendly neighbourhood Spider-Man"
git config --local user.email "re-buildrun-team@digital.cabinet-office.gov.uk"
git remote add source "${SOURCE_REPO_URL}"
git remote add destination "${CODECOMMIT_REPO_URL}"
git fetch --all
git pull source master
dest_master="$(git ls-remote destination master)"

if ! [ -z "$dest_master" ]
then
  echo "Destination repository $CODECOMMIT_REPO_URL already has master branch, nothing to do"
  exit 0
fi

git branch -u source/master

# update the embedded timestamp
sed -i.bak -E -e "s/(chartCommitTimestamp: )\"[0-9]+\"/\1\"$(date +%s)\"/g" gsp-canary/charts/values.yaml
git add gsp-canary/charts/values.yaml
git commit -m "Initial timestamp update."
git push --force destination master
