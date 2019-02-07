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
