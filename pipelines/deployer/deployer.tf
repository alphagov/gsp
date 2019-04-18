terraform {
  backend "s3" {}
}

variable "aws_account_role_arn" {
  type = "string"
}

variable "account_name" {
  type = "string"
}

variable "cluster_name" {
  type = "string"
}

variable "cluster_domain" {
  type = "string"
}

variable "github_client_id" {
  type = "string"
}

variable "github_client_secret" {
  type = "string"
}

variable "splunk_enabled" {
  type    = "string"
  default = "0"
}

variable "splunk_hec_url" {
  type = "string"
}

variable "splunk_hec_token" {
  type = "string"
}

variable "splunk_index" {
  type    = "string"
  default = "run_sandbox_k8s"
}

variable "worker_instance_type" {
  type    = "string"
  default = "m5.large"
}

variable "worker_count" {
  type    = "string"
  default = "3"
}

variable "ci_worker_instance_type" {
  type    = "string"
  default = "m5.large"
}

variable "ci_worker_count" {
  type    = "string"
  default = "3"
}

variable "eks_version" {
  description = "EKS platform version (https://docs.aws.amazon.com/eks/latest/userguide/platform-versions.html)"
  type        = "string"
}

provider "aws" {
  region = "eu-west-2"

  assume_role {
    role_arn = "${var.aws_account_role_arn}"
  }
}

data "aws_caller_identity" "current" {}

module "gsp-domain" {
  source         = "../../modules/gsp-domain"
  existing_zone  = "govsvc.uk"
  delegated_zone = "${var.cluster_domain}"

  providers = {
    aws = "aws"
  }
}

module "gsp-network" {
  source       = "../../modules/gsp-network"
  cluster_name = "${var.cluster_name}"
}

module "gsp-cluster" {
  source            = "../../modules/gsp-cluster"
  account_name      = "${var.account_name}"
  cluster_name      = "${var.cluster_name}"
  cluster_domain    = "${var.cluster_domain}"
  cluster_domain_id = "${module.gsp-domain.zone_id}"

  admin_role_arns = [
    "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/admin",
    "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/deployer",
  ]

  admin_user_arns = [
    "arn:aws:iam::622626885786:user/chris.farmiloe@digital.cabinet-office.gov.uk",
  ]

  sre_user_arns = [
    "arn:aws:iam::622626885786:user/chris.farmiloe@digital.cabinet-office.gov.uk",
  ]

  gds_external_cidrs = [
    "213.86.153.212/32",
    "213.86.153.213/32",
    "213.86.153.214/32",
    "213.86.153.235/32",
    "213.86.153.236/32",
    "213.86.153.237/32",
    "85.133.67.244/32",
    "18.130.144.30/32",  # autom8 concourse
    "3.8.110.67/32",     # autom8 concourse
  ]

  eks_version             = "${var.eks_version}"
  worker_instance_type    = "${var.worker_instance_type}"
  worker_count            = "${var.worker_count}"
  ci_worker_instance_type = "${var.ci_worker_instance_type}"
  ci_worker_count         = "${var.ci_worker_count}"

  vpc_id                 = "${module.gsp-network.vpc_id}"
  private_subnet_ids     = "${module.gsp-network.private_subnet_ids}"
  public_subnet_ids      = "${module.gsp-network.public_subnet_ids}"
  nat_gateway_public_ips = "${module.gsp-network.nat_gateway_public_ips}"
  splunk_hec_url         = "${var.splunk_hec_url}"
  splunk_hec_token       = "${var.splunk_hec_token}"
  splunk_index           = "${var.splunk_index}"

  codecommit_init_role_arn = "${var.aws_account_role_arn}"
  github_client_id         = "${var.github_client_id}"
  github_client_secret     = "${var.github_client_secret}"
}

# module "prototype-kit" {
#   source = "git::https://github.com/alphagov/gsp-terraform-ignition//modules/flux-release"

#   namespace      = "gsp-prototype-kit"
#   chart_git      = "https://github.com/alphagov/gsp-govuk-prototype-kit.git"
#   chart_ref      = "gsp"
#   chart_path     = "charts/govuk-prototype-kit"
#   cluster_name   = "${var.cluster_domain}"
#   cluster_domain = "${var.cluster_domain}"
#   addons_dir     = "addons/${var.cluster_domain}"

#   values = <<EOF
#     ingress:
#       hosts:
#         - pk.${var.cluster_domain}
#         - prototype-kit.${var.cluster_domain}
#       tls:
#         - secretName: prototype-kit-tls
#           hosts:
#             - pk.${var.cluster_domain}
#             - prototype-kit.${var.cluster_domain}
# EOF
# }

output "kubeconfig" {
  value = "${module.gsp-cluster.kubeconfig}"
}

output "values" {
  sensitive = true
  value     = "${module.gsp-cluster.values}"
}
