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
sed -i "" -E -e "s/(chartCommitTimestamp: )\"[0-9]+\"/\1\"$(date +%s)\"/g" charts/gsp-canary/values.yaml
git add charts/gsp-canary/values.yaml
git commit -m "Initial timestamp update."
git push --force destination master
