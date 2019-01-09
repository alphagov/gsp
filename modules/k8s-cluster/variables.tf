variable "cluster_domain_suffix" {
  type = "string"
}

variable "kubelet_kubeconfig" {
  type = "string"
}

variable "kube_ca_crt" {
  type = "string"
}

variable "user_data_bucket_name" {
  type = "string"
}

variable "user_data_bucket_region" {
  type = "string"
}

variable "vpc_id" {
  type = "string"
}

variable "subnet_ids" {
  type = "list"
}

variable "controller_target_group_arns" {
  type = "list"
}

variable "worker_target_group_arns" {
  type = "list"
}

variable "cluster_name" {
  type = "string"
}

variable "k8s_tag" {
  type = "string"
}

variable "s3_user_data_policy_arn" {
  type = "string"
}

variable "service_cidr" {
  type    = "string"
  default = "10.3.0.0/24"
}

variable "controller_node_labels" {
  type    = "string"
  default = "node-role.kubernetes.io/master"
}

variable "controller_node_taints" {
  type    = "string"
  default = "node-role.kubernetes.io/master=:NoSchedule"
}

variable "worker_node_labels" {
  type    = "string"
  default = "node-role.kubernetes.io/node"
}

variable "worker_node_taints" {
  type    = "string"
  default = ""
}

variable "controller_instance_type" {
  type    = "string"
  default = "t2.small"
}

variable "worker_instance_type" {
  type    = "string"
  default = "t2.small"
}

variable "controller_count" {
  type    = "string"
  default = "1"
}

variable "worker_count" {
  type    = "string"
  default = "2"
}
