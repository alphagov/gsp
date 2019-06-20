variable "account_name" {
  type        = "string"
  default     = "gds"
  description = "descriptive label of account, programme department who owns this cluster"
}

variable "cluster_name" {
  type = "string"
}

variable "cluster_domain" {
  description = "The FQDN of the DNS zone allocated to this cluster"
  type        = "string"
}

variable "cluster_domain_id" {
  description = "The zone id of DNS zone allocated to this cluster"
  type        = "string"
}

variable "admin_role_arns" {
  type = "list"
}

variable "eks_version" {
  type = "string"
}

variable "dev_user_arns" {
  description = "A list of user ARNs that will be mapped to the cluster dev role"
  type        = "list"
  default     = []
}

variable "admin_user_arns" {
  description = "A list of user ARNs that will be mapped to the cluster admin role"
  type        = "list"
  default     = []
}

variable "sre_role_arns" {
  description = "A list of role ARNs that will be mapped to the cluster sre role"
  type        = "list"
  default     = []
}

variable "sre_user_arns" {
  description = "A list of user ARNs that will be mapped to the cluster sre role"
  type        = "list"
  default     = []
}

variable "worker_count" {
  type    = "string"
  default = "2"
}

variable "worker_instance_type" {
  type    = "string"
  default = "t2.small"
}

variable "ci_worker_count" {
  type    = "string"
  default = "2"
}

variable "ci_worker_instance_type" {
  type    = "string"
  default = "m5d.large"
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

variable "splunk_enabled" {
  description = "Enable splunk log shipping"
  type        = "string"
  default     = "0"
}

variable "splunk_hec_url" {
  description = "Splunk HTTP event collector URL to send logs to"
  type        = "string"
  default     = ""
}

variable "k8s_splunk_hec_token" {
  description = "Splunk HTTP event collector token for authentication"
  type        = "string"
  default     = ""
}

variable "k8s_splunk_index" {
  description = "Name of index to be added as metadata to logs for use in splunk"
  type        = "string"
  default     = ""
}

variable "vpc_flow_log_splunk_hec_token" {
  description = "Splunk HTTP event collector token for authentication"
  type        = "string"
  default     = ""
}

variable "vpc_flow_log_splunk_index" {
  description = "Name of index to be added as metadata to logs for use in splunk"
  type        = "string"
  default     = ""
}

variable "cloud_hsm_ip" {
  type = "string"
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

variable "egress_ips" {
  type = "list"
}

variable "ingress_ips" {
  type = "list"
}

variable "github_teams" {
  default     = ["alphagov:re-gsp"]
  description = "the list of github teams allowed to be authenticated into concourse"
}

variable "github_client_id" {
  description = "the github application client_id ID to allow oauth"
}

variable "github_client_secret" {
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

variable "concourse_main_team_github_teams" {
  default     = ["alphagov:re-gsp"]
  description = "the list of github teams authorized to view the concourse 'main' team"
}
