variable "cluster_name" {
  type = "string"
}

variable "private_cidr_block" {
  description = "CIDR IPv4 range for private subnet"
  type        = "string"
}

variable "public_cidr_block" {
  description = "CIDR IPv4 range for public subnet"
  type        = "string"
}

variable "vpc_id" {
  description = "VPC ID"
}

variable "availability_zone" {
  description = "The availability zone for this subnet"
}

variable "internet_gateway_id" {
  description = "gateway id for public subnet"
}
