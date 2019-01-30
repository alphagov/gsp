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
