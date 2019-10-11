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

variable "worker_eks_version" {
  type = "string"
}

variable "worker_instance_type" {
  type    = "string"
  default = "t3.medium"
}

variable "minimum_workers_per_az_count" {
  type    = "string"
  default = "1"
}

variable "maximum_workers_per_az_count" {
  type    = "string"
  default = "5"
}

variable "ci_worker_instance_type" {
  type    = "string"
  default = "t3.medium"
}

variable "ci_worker_count" {
  type    = "string"
  default = "3"
}
