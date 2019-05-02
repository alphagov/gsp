#!/bin/bash

set -eu

# rough check for missing vars in terraform
(cd pipelines/deployer \
	&& terraform init --backend=false \
	&& terraform validate \
		--var account_name=x \
		--var splunk_enabled=0 \
		--var splunk_hec_token=x \
		--var splunk_hec_url=x \
		--var github_client_secret=x \
		--var github_client_id=x \
		--var cluster_name=x \
		--var cluster_domain=x \
		--var aws_account_role_arn=x \
		--var eks_version=x \
)


# rough check for missing vars in values.yaml for gsp-cluster chart
rm -rf output || echo 'ok'
mkdir -p output

cat << 'EOF' > output/values.yaml
namespaces:
- name: sandbox-canary
  resources: [{"apiVersion":"v1","kind":"Secret","type":"Opaque","metadata":{"name":"ci-deploy-key"},"data":{"private_key":"RklYTUUK"}},{"apiVersion":"v1","kind":"ConfigMap","metadata":{"name":"ci-deploy-key"},"data":{"public_key":"RklYTUUK"}},{"apiVersion":"flux.weave.works/v1beta1","kind":"HelmRelease","metadata":{"name":"canary"},"spec":{"releaseName":"canary","chart":{"git":"{{ .Values.global.canary.repository }}","ref":"master","path":"charts/gsp-canary","verificationKeys":"{{ .Values.global.canary.verificationKeys }}"},"values":{"annotations":{"iam.amazonaws.com/role":"{{ .Values.global.roles.canary }}"},"updater":{"helmChartRepoUrl":"{{ .Values.global.canary.repository }}"}}}}]
users:
- name: chris.farmiloe
  email: chris.farmiloe@digital.cabinet-office.gov.uk
  ARN: arn:aws:iam::622626885786:user/chris.farmiloe@digital.cabinet-office.gov.uk
  roles:
  - sandbox-admin
  - sandbox-sre
  - sandbox-canary-dev
  roleARN: arn:aws:iam::011571571136:role/farms-user-chris.farmiloe
- name: sam.crang
  email: sam.crang@digital.cabinet-office.gov.uk
  ARN: arn:aws:iam::622626885786:user/sam.crang@digital.cabinet-office.gov.uk
  roles:
  - sandbox-canary-dev
  roleARN: arn:aws:iam::011571571136:role/farms-user-sam.crang
EOF

helm template \
	--output-dir ./output \
	--namespace fake-system \
	--values output/values.yaml \
	--values <(\
		cat modules/gsp-cluster/data/values.yaml \
		| sed 's/${admin_role_arns}/[]/' \
		| sed 's/${admin_user_arns}/[]/' \
		| sed 's/${sre_role_arns}/[]/' \
		| sed 's/${sre_user_arns}/[]/' \
		| sed 's/${bootstrap_role_arns}/[]/' \
		| sed 's/${concourse_teams}/["org:team"]/' \
	) \
	"charts/gsp-cluster"

