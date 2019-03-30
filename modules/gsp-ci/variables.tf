variable "enabled" {
  default = 1
}

variable "cluster_name" {
  type = "string"
}

variable "dns_zone" {
  type = "string"
}

variable "harbor_role_asumer_arn" {
  type = "string"
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
