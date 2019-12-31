#!/bin/bash

set -eu


# while there is a transition to helm v3 you may need to have two versions of helm
# there's no such thing as helmenv or hvm so having two binaries is best bet for now
# we want helm2, so if there's a helm2, use that
helm="helm"
if [ -x "$(command -v helm2)" ]; then
	helm="helm2"
fi

# we don't build the lifecycle hook here, but it will complain without it
touch modules/k8s-cluster/aws-node-lifecycle-hook.zip

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
  account:
    name: sandbox
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
- name: test-operators
external_dns:
- namespace: gsp-system
  role_name: sandbox-gsp-system-external-dns
  zone_id: gspsystem
- namespace: verify-proxy-node-build
  role_name: sandbox-verify-proxy-node-build-external-dns
  zone_id: proxynode
users:
- name: chris.farmiloe
  email: chris.farmiloe@digital.cabinet-office.gov.uk
  ARN: arn:aws:iam::000000072:user/chris.farmiloe@digital.cabinet-office.gov.uk
  roles:
  - role: admin
    account: sandbox
  roleARN: arn:aws:iam::000000072:role/chris.farmiloe
  github: "chrisfarms"
- name: sam.crang
  email: sam.crang@digital.cabinet-office.gov.uk
  ARN: arn:aws:iam::000000072:user/sam.crang@digital.cabinet-office.gov.uk
  roles:
  - role: operator
    namespace: test-operators
    account: sandbox
  roleARN: arn:aws:iam::000000072:role/sam.crang
  github: "samcrang"
- name: jeff.jefferson
  email: jeff.jefferson@digital.cabinet-office.gov.uk
  ARN: arn:aws:iam::000000072:user/jeff.jefferson@digital.cabinet-office.gov.uk
  roles:
  - role: dev
    account: portfolio
  roleARN: arn:aws:iam::000000072:role/jeff.jefferson
  github: "jefferz83"
extraPermissionsDev: []
extraPermissionsSRE: []
EOF

gomplate -d config=output/values.yaml -f templates/managed-namespaces-gateways.yaml > output/gateways-values.yaml

gomplate -d config=output/values.yaml -f templates/managed-namespaces-zones.tf > pipelines/deployer/managed-namespaces-zones.tf

# rough check for missing vars in terraform
(cd pipelines/deployer \
	&& terraform init --backend=false \
	&& terraform validate)

$helm template \
	--output-dir ./output \
	--name gsp \
	--namespace gsp-system \
	--values <(\
		cat modules/gsp-cluster/data/values.yaml \
		| sed 's/${admin_role_arns}/[]/' \
		| sed 's/${bootstrap_role_arns}/[]/' \
		| sed 's/${concourse_teams}/["org:team"]/' \
		| sed 's/${egress_ip_addresses}/[]/' \
		| sed 's/${eks_version}/1.14/' \
		| sed 's/${external_dns_map}/external_dns: []/' \
	) \
	--values output/values.yaml \
	--set 'global.cloudHsm.enabled=true' \
	"charts/gsp-cluster"

$helm template \
--name istio \
--namespace istio-system \
--output-dir output \
--set global.runningOnAws=true \
--values output/gateways-values.yaml \
charts/gsp-istio
