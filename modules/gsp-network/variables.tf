variable "cluster_name" {
  type = string
}

variable "netnum" {
  description = "network number (0-255) for assigned 10.x.0.0/16 cidr, preferably unique per persistant cluster"
  default     = "0"
}

