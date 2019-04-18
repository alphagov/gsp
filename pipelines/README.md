# Applying pipeline

## Before you start

You need the following dependencies installed on your laptop: [`yq`](https://pypi.org/project/yq/), `jq`.

## How to update the old verify pipeline

```
fly -t gsp set-pipeline -p $ACCOUNT_NAME \
	--config tools-staging-prod-infra.yaml \
	--var account-name=$ACCOUNT_NAME \
	--var account-role-arn=$DEPLOYER_ROLE_ARN \
	--yaml-var public-gpg-keys="$(yq . ../users/*.yaml | jq -s '[.[] | select(.teams[] | IN("re-gsp")) | .pub]')" \
	--check-creds
```

## How to create or update a deployer pipeline

This example will push a pipeline for a cluster in the "sandbox" account called "farms" with a pipeline called "farms".

```
fly -t deployer set-pipeline -p farms --config pipelines/deployer/deployer.yaml \
  --var account-name=sandbox \
  --var account-role-arn=arn:aws:iam::011571571136:role/deployer \
  --var cluster-name=farms \
  --yaml-var trusted-developer-keys="[]" \
  --var splunk_hec_token="NOTATOKEN" \
  --var github-client-secret=NOTASECRET \
  --var github-client-id=NOTID \
  --var eks-version=1.12 \
  --var splunk-enabled=0 \
  --var splunk-hec-token=NOTOKEN \
  --var splunk-hec-url=whataver \
  --check-creds
```
