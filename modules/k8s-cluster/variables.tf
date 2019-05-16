variable "vpc_id" {
  type = "string"
}

variable "public_subnet_ids" {
  type = "list"
}

variable "private_subnet_ids" {
  type = "list"
}

variable "cluster_name" {
  type = "string"
}

variable "apiserver_allowed_cidrs" {
  type = "list"
}

variable "eks_version" {
  type = "string"
}

variable "worker_instance_type" {
  type    = "string"
  default = "t2.small"
}

variable "worker_count" {
  type    = "string"
  default = "3"
}

variable "ci_worker_instance_type" {
  type    = "string"
  default = "t2.small"
}

variable "ci_worker_count" {
  type    = "string"
  default = "3"
}
