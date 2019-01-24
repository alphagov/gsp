variable "apiserver_address" {
  type = "string"
}

variable "cluster_domain_suffix" {
  type = "string"
}

variable "etcd_servers" {
  type = "list"
}

variable "k8s_tag" {
  type = "string"
}

variable "cluster_name" {
  type = "string"
}

variable "cluster_id" {
  type = "string"
}

variable "etcd_ca_cert_pem" {
  type = "string"
}

variable "etcd_client_private_key_pem" {
  type = "string"
}

variable "etcd_client_cert_pem" {
  type = "string"
}

variable "admin_role_arns" {
  description = "A list of ARNs that will be mapped to cluster administrators"
  type        = "list"
}

variable "admin_role_arn_mapping_template" {
  description = "The template that renders into yaml for the aws iam authenticator. Whitespace is important here."
  type        = "string"

  default = <<TEMPLATE
      - roleARN: %s
        username: admin
        groups:
        - system:masters
TEMPLATE
}

variable "dev_role_arns" {
  description = "A list of ARNs that will be mapped to cluster devs"
  type        = "list"
  default     = []
}

variable "dev_role_arn_mapping_template" {
  description = "The template that renders into yaml for the aws iam authenticator. Whitespace is important here."
  type        = "string"

  default = <<TEMPLATE
      - roleARN: %s
        username: dev
        groups:
        - dev
TEMPLATE
}

variable "assets_dir" {
  type    = "string"
  default = "/opt/bootkube/assets"
}

variable "etcd_port" {
  type    = "string"
  default = "2379"
}

variable "service_cidr" {
  type    = "string"
  default = "10.3.0.0/24"
}

variable "pod_cidr" {
  type    = "string"
  default = "10.2.0.0/16"
}
