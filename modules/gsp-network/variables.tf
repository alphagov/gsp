variable "cluster_name" {
  type = "string"
}

variable "cidr_block" {
  description = "CIDR IPv4 range of the VPC"
  type        = "string"
  default     = "10.0.0.0/16"
}
