variable "cluster_name" {
  type = "string"
}

variable "dns_zone" {
  type = "string"
}

variable "subnet_ids" {
  type = "list"
}

variable "vpc_id" {
  type = "string"
}

variable "dns_zone_id" {
  type = "string"
}

variable "user_data_bucket_name" {
  type = "string"
}

variable "s3_user_data_policy_arn" {
  type = "string"
}

variable "node_count" {
  type    = "string"
  default = "3"
}

variable "instance_type" {
  type    = "string"
  default = "t2.small"
}

variable "assets_dir" {
  type    = "string"
  default = "/opt/bootkube/assets"
}
