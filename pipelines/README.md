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

This example will push a pipeline for a cluster in the "sandbox" account called "gsp"

```
fly -t cd-gsp set-pipeline -p sandbox --config pipelines/deployer.yaml --var account-name=sandbox --var account-role-arn=arn:aws:iam::011571571136:role/deployer --var cluster-name=gsp --yaml-var trusted-developer-keys="$(yq . ./users/*.yaml | jq -s '[ .[].pub ]')" --var splunk_hec_token="NOTATOKEN" --var github-client-secret=NOTASECRET --var github-client-id=NOTID --check-creds
```
