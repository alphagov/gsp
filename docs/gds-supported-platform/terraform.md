# GDS Supported Platform Terraform
Terraform for the various parts of the GDS Supported Platform

## Intention

A collection of terraform modules that comprise the fundamental components of a GSP cluster. AWS EKS provides the kubernetes foundation, with several modules split out separately for "persistence" so the cluster can be torn down and rebuilt without losing some longer-term values (such as sealed secrets keys).

## Usage

The main entrypoint of this repo is the `gsp-cluster` terraform module. An example might look like:

```
data "aws_caller_identity" "current" {}

module "gsp-network" {
  source       = "git::https://github.com/alphagov/gsp-terraform-ignition//modules/gsp-network"
  cluster_name = "rafalp"
}

module "gsp-persistent" {
  source       = "git::https://github.com/alphagov/gsp-terraform-ignition//modules/gsp-persistent"
  cluster_name = "${module.gsp-network.cluster-name}"
  cluster_domain     = "run-sandbox.aws.ext.govsvc.uk"
}

module "gsp-cluster" {
    source = "git::https://github.com/alphagov/gsp-terraform-ignition//modules/gsp-cluster"
    cluster_name = "rafalp"
    cluster_domain = "run-sandbox.aws.ext.govsandbox.uk"
    admin_role_arns = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/admin"]
    worker_instance_type = "m5.large"
    worker_count = "2"

    addons = {
      ci = 1
    }

    sealed_secrets_cert_pem        = "${module.gsp-persistent.sealed_secrets_cert_pem}"
    sealed_secrets_private_key_pem = "${module.gsp-persistent.sealed_secrets_private_key_pem}"
    vpc_id                         = "${module.gsp-network.vpc_id}"
    private_subnet_ids             = "${module.gsp-network.private_subnet_ids}"
    public_subnet_ids              = "${module.gsp-network.public_subnet_ids}"
    nat_gateway_public_ips         = "${module.gsp-network.nat_gateway_public_ips}"

    sre_user_arns = ["arn:aws:iam::622626885786:user/rafal.proszowski@digital.cabinet-office.gov.uk"]

    github_client_id = "1234567890"
    github_client_secret = "qwertyuiop"
}

output "kubeconfig" {
    value = "${module.gsp-cluster.kubeconfig}"
}

```
