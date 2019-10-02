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
  default = "t3.medium"
}

variable "worker_count" {
  type    = "string"
  default = "3"
}

variable "extra_workers_per_az_count" {
  type    = "string"
  default = "0"
}

variable "ci_worker_instance_type" {
  type    = "string"
  default = "t3.medium"
}

variable "ci_worker_count" {
  type    = "string"
  default = "3"
}
