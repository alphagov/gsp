variable "bootstrap_base_userdata_source" {
  type = "string"
}

variable "bootstrap_base_userdata_verification" {
  type = "string"
}

variable "user_data_bucket_name" {
  type = "string"
}

variable "user_data_bucket_region" {
  type = "string"
}

variable "cluster_name" {
  type = "string"
}

variable "security_group_ids" {
  type = "list"
}

variable "subnet_id" {
  type = "string"
}

variable "iam_instance_profile_name" {
  type = "string"
}

variable "lb_target_group_arn" {
  type = "string"
}

variable "dns_service_ip" {
  type = "string"
}

variable "cluster_domain_suffix" {
  type = "string"
}

variable "k8s_tag" {
  type = "string"
}

variable "kubelet_kubeconfig" {
  type = "string"
}

variable "kube_ca_crt" {
  type = "string"
}

variable "instance_type" {
  type    = "string"
  default = "t2.small"
}
