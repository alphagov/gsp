variable "cluster_name" {
  type = "string"
}

variable "dns_zone" {
  type = "string"
}

variable "addons_dir" {
  description = "local target path to place kubernetes resource yaml"
  type        = "string"
  default     = "addons"
}

variable "canary_role_assumer_arn" {
  description = "ARN of the role assuming the canary role, e.g. kiam-server-role"
  type        = "string"
}
