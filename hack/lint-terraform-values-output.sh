#!/bin/bash

set -eu

# rough check for missing vars in terraform
(cd pipelines/deployer \
	&& terraform init --backend=false \
	&& terraform validate \
		--var account_name=x \
		--var splunk_enabled=0 \
		--var splunk_hec_url=x \
		--var k8s_splunk_hec_token=x \
		--var k8s_splunk_index=x \
		--var hsm_splunk_hec_token=x \
		--var hsm_splunk_index=x \
		--var vpc_flow_log_splunk_hec_token=x \
		--var vpc_flow_log_splunk_index=x \
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
global:
  cluster:
    domain: fake.com
    name: sandbox
    privateKey: BEGIN-PRIVATE
    publicKey: BEGIN-PUBLIC
namespaces:
- name: verify-metadata-controller
  owner: alphagov
  repository: verify-metadata-controller
  branch: master
  path: ci
  permittedRolesRegex: "^$"
  requiredApprovalCount: 2
  scope: cluster
- name: verify-proxy-node-build
  owner: alphagov
  repository: verify-proxy-node
  path: ci/build
  requiredApprovalCount: 2
users:
- name: chris.farmiloe
  email: chris.farmiloe@digital.cabinet-office.gov.uk
  ARN: arn:aws:iam::000000072:user/chris.farmiloe@digital.cabinet-office.gov.uk
  roles:
  - role: admin
    account: sandbox
  roleARN: arn:aws:iam::000000072:role/chris.farmiloe
  github: "chrisfarms"
  pub: BEGIN-FARMS-KEY
- name: sam.crang
  email: sam.crang@digital.cabinet-office.gov.uk
  ARN: arn:aws:iam::000000072:user/sam.crang@digital.cabinet-office.gov.uk
  roles:
  - role: sre
    account: sandbox
  roleARN: arn:aws:iam::000000072:role/sam.crang
  github: "samcrang"
  pub: BEGIN-SAM-KEY
- name: jeff.jefferson
  email: jeff.jefferson@digital.cabinet-office.gov.uk
  ARN: arn:aws:iam::000000072:user/jeff.jefferson@digital.cabinet-office.gov.uk
  roles:
  - role: dev
    account: portfolio
  roleARN: arn:aws:iam::000000072:role/jeff.jefferson
  github: "jefferz83"
  pub: BEGIN-JEFF-KEY
EOF

helm template \
	--output-dir ./output \
	--name gsp \
	--namespace gsp-system \
	--values <(\
		cat modules/gsp-cluster/data/values.yaml \
		| sed 's/${admin_role_arns}/[]/' \
		| sed 's/${admin_user_arns}/[]/' \
		| sed 's/${sre_role_arns}/[]/' \
		| sed 's/${sre_user_arns}/[]/' \
		| sed 's/${bootstrap_role_arns}/[]/' \
		| sed 's/${concourse_teams}/["org:team"]/' \
	) \
	--values output/values.yaml \
	"charts/gsp-cluster"

