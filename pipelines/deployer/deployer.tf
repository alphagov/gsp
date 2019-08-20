variable "account_name" {
  type = "string"
}

variable "cluster_name" {
  type = "string"
}

variable "cluster_number" {
  description = "unique number (0-255) for this cluster, dictates the assigned network 10.x.0.0/16"
  default     = "0"
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

variable "k8s_splunk_hec_token" {
  type = "string"
}

variable "k8s_splunk_index" {
  type = "string"
}

variable "hsm_splunk_hec_token" {
  type = "string"
}

variable "hsm_splunk_index" {
  type = "string"
}

variable "vpc_flow_log_splunk_hec_token" {
  type = "string"
}

variable "vpc_flow_log_splunk_index" {
  type = "string"
}

variable "worker_instance_type" {
  type    = "string"
  default = "m5d.large"
}

variable "worker_count" {
  type    = "string"
  default = "3"
}

variable "ci_worker_instance_type" {
  type    = "string"
  default = "m5d.large"
}

variable "ci_worker_count" {
  type    = "string"
  default = "3"
}

variable "eks_version" {
  description = "EKS platform version (https://docs.aws.amazon.com/eks/latest/userguide/platform-versions.html)"
  type        = "string"
}

variable "enable_nlb" {
  type    = "string"
  default = "0"
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
  netnum       = "${var.cluster_number}"
}

module "hsm" {
  source                   = "../../modules/hsm"
  subnet_cidr_map          = "${module.gsp-network.private_subnet_cidr_mapping}"
  source_security_group_id = "${module.gsp-cluster.worker_security_group_id}"
  cluster_name             = "${var.cluster_name}"
  splunk                   = "${var.splunk_enabled}"
  splunk_hec_url           = "${var.splunk_hec_url}"
  splunk_hec_token         = "${var.hsm_splunk_hec_token}"
  splunk_index             = "${var.hsm_splunk_index}"
}

module "gsp-cluster" {
  source            = "../../modules/gsp-cluster"
  account_name      = "${var.account_name}"
  cluster_name      = "${var.cluster_name}"
  cluster_domain    = "${var.cluster_domain}"
  cluster_domain_id = "${module.gsp-domain.zone_id}"

  admin_role_arns = [
    "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/deployer",
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

  vpc_id                        = "${module.gsp-network.vpc_id}"
  private_subnet_ids            = "${module.gsp-network.private_subnet_ids}"
  public_subnet_ids             = "${module.gsp-network.public_subnet_ids}"
  egress_ips                    = "${module.gsp-network.egress_ips}"
  ingress_ips                   = "${module.gsp-network.ingress_ips}"
  splunk_enabled                = "${var.splunk_enabled}"
  splunk_hec_url                = "${var.splunk_hec_url}"
  k8s_splunk_hec_token          = "${var.k8s_splunk_hec_token}"
  k8s_splunk_index              = "${var.k8s_splunk_index}"
  vpc_flow_log_splunk_hec_token = "${var.vpc_flow_log_splunk_hec_token}"
  vpc_flow_log_splunk_index     = "${var.vpc_flow_log_splunk_index}"

  github_client_id     = "${var.github_client_id}"
  github_client_secret = "${var.github_client_secret}"

  enable_nlb = "${var.enable_nlb}"
}

output "kubeconfig" {
  value = "${module.gsp-cluster.kubeconfig}"
}

output "values" {
  sensitive = true
  value     = "${module.gsp-cluster.values}"
}

output "gsp_istio_values" {
  sensitive = true
  value = "${module.gsp-cluster.gsp_istio_values}"
}

output "vpc_id" {
  value = "${module.gsp-network.vpc_id}"
}

output "subnet_ids" {
  value = "${module.gsp-network.private_subnet_ids}"
}

output "worker_security_group_id" {
  value = "${module.gsp-cluster.worker_security_group_id}"
}

output "r53_zone_id" {
  value = "${module.gsp-domain.zone_id}"
}

output "r53_zone_name" {
  value = "${module.gsp-domain.name}"
}

output "cluster_domain" {
  value = "${var.cluster_domain}"
}
