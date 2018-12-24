variable "cluster_name" {
  type = "string"
}

variable "cluster_id" {
  type = "string"
}

variable "dns_zone_id" {
  type = "string"
}

variable "dns_zone" {
  type = "string"
}

variable "user_data_bucket_name" {
  type = "string"
}

variable "user_data_bucket_region" {
  type = "string"
}

variable "admin_role_arns" {
  type = "list"
}

variable "host_cidr" {
  description = "CIDR IPv4 range to assign to EC2 nodes"
  type        = "string"
  default     = "10.0.0.0/16"
}

variable "etcd_node_count" {
  type    = "string"
  default = "3"
}

variable "k8s_tag" {
  type    = "string"
  default = "v1.12.2"
}

variable "controller_count" {
  type    = "string"
  default = "1"
}

variable "worker_count" {
  type    = "string"
  default = "2"
}

variable "etcd_instance_type" {
  type    = "string"
  default = "t2.small"
}

variable "controller_instance_type" {
  type    = "string"
  default = "t2.small"
}

variable "worker_instance_type" {
  type    = "string"
  default = "t2.small"
}
