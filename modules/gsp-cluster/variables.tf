variable "cluster_name" {
  type = "string"
}

variable "dns_zone" {
  type = "string"
}

variable "user_data_bucket_name" {
  type = "string"
}

variable "user_data_bucket_region" {
  type = "string"
}

variable "admin_role_arns" {
  type = "list"
}

variable "dev_user_arns" {
  description = "A list of user ARNs that will be mapped to the cluster dev role"
  type        = "list"
  default     = []
}

variable "sre_user_arns" {
  description = "A list of user ARNs that will be mapped to the cluster sre role"
  type        = "list"
  default     = []
}

variable "host_cidr" {
  description = "CIDR IPv4 range to assign to EC2 nodes"
  type        = "string"
}

variable "etcd_node_count" {
  type    = "string"
  default = "3"
}

variable "k8s_tag" {
  type    = "string"
  default = "v1.12.2"
}

variable "controller_count" {
  type    = "string"
  default = "1"
}

variable "worker_count" {
  type    = "string"
  default = "2"
}

variable "etcd_instance_type" {
  type    = "string"
  default = "t2.small"
}

variable "controller_instance_type" {
  type    = "string"
  default = "t2.small"
}

variable "worker_instance_type" {
  type    = "string"
  default = "t2.small"
}

variable "addons" {
  type = "map"

  default = {}
}

variable "gds_external_cidrs" {
  description = "External GDS CIDRs that are allowed to talk to the clusters, taken from the GDS wiki"
  type        = "list"

  default = [
    "213.86.153.212/32",
    "213.86.153.213/32",
    "213.86.153.214/32",
    "213.86.153.235/32",
    "213.86.153.236/32",
    "213.86.153.237/32",
    "85.133.67.244/32",
  ]
}

variable "dev_namespaces" {
  type    = "list"
  default = []
}

variable "splunk_hec_token" {
  description = "Splunk HTTP event collector token for authentication"
  type        = "string"
  default     = ""
}

variable "splunk_hec_url" {
  description = "Splunk HTTP event collector URL to send logs to"
  type        = "string"
  default     = ""
}

variable "splunk_index" {
  description = "Name of index to be added as metadata to logs for use in splunk"
  type        = "string"
  default     = ""
}

variable "codecommit_init_role_arn" {
  type    = "string"
  default = ""
}

variable "vpc_id" {
  type = "string"
}

variable "private_subnet_ids" {
  type = "list"
}

variable "public_subnet_ids" {
  type = "list"
}

variable "nat_gateway_public_ips" {
  type = "list"
}

variable "sealed_secrets_cert_pem" {
  description = "Sealed secrets cert"
  type        = "string"
}

variable "sealed_secrets_private_key_pem" {
  description = "Sealed secrets private key"
  type        = "string"
}

variable "github_teams" {
  default     = ["alphagov:re-gsp"]
  description = "the list of github teams allowed to be authenticated into concourse"
}

variable "github_client_id" {
  default     = ""
  description = "the github application client_id ID to allow oauth"
}

variable "github_client_secret" {
  default     = ""
  description = "the github application client_secret ID to allow oauth"
}

variable "github_ca_cert" {
  default     = ""
  description = "the github application ca_cert ID to allow oauth"
}

variable "concourse_teams" {
  default     = []
  description = "the list of teams to be created in concourse"
}
