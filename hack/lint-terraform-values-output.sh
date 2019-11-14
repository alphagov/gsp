#!/bin/bash

set -eu

touch modules/k8s-cluster/aws-node-lifecycle-hook.zip

# rough check for missing vars in terraform
(cd pipelines/deployer \
	&& terraform init --backend=false \
	&& terraform validate)


# rough check for missing vars in values.yaml for gsp-cluster chart
rm -rf output || echo 'ok'
mkdir -p output

cat << 'EOF' > output/values.yaml
global:
  cluster:
    domain: fake.com
    name: sandbox
    egressIpAddresses: ["127.0.0.1", "1.2.3.4"]
    privateKey: BEGIN-PRIVATE
    publicKey: BEGIN-PUBLIC
egressSafelist:
- name: integration-hub
  service:
    hosts: ["www.integration.signin.service.gov.uk"]
    ports:
    - name: verify-integration-https
      number: 443
      protocol: TCP
    location: MESH_EXTERNAL
    resolution: DNS
- name: production-hub
  service:
    hosts: ["www.signin.service.gov.uk"]
    ports:
    - name: verify-https
      number: 443
      protocol: TCP
    location: MESH_EXTERNAL
    resolution: DNS
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
  ingress:
    enabled: true
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
extraPermissionsDev: []
extraPermissionsSRE: []
EOF

gomplate -d config=output/values.yaml -f templates/managed-namespaces-gateways.yaml > output/gateways-values.yaml

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
		| sed 's/${egress_ip_addresses}/[]/' \
		| sed 's/${eks_version}/1.14/' \
	) \
	--values output/values.yaml \
	--set 'global.cloudHsm.enabled=true' \
	"charts/gsp-cluster"

helm template \
--name istio \
--namespace istio-system \
--output-dir output \
--set global.runningOnAws=true \
--values output/gateways-values.yaml \
charts/gsp-istio
