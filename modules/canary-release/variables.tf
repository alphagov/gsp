variable "cluster_id" {
  type = "string"
}

variable "addons_dir" {
  description = "local target path to place kubernetes resource yaml"
  type        = "string"
  default     = "addons"
}
