variable "vpc_id" {
  type = "string"
}

variable "subnet_ids" {
  type = "list"
}

variable "cluster_name" {
  type = "string"
}

variable "apiserver_allowed_cidrs" {
  type = "list"
}

variable "worker_instance_type" {
  type    = "string"
  default = "t2.small"
}

variable "worker_count" {
  type    = "string"
  default = "2"
}

variable "admin_role_arns" {
  description = "A list of ARNs that will be mapped to cluster administrators"
  type        = "list"
}

variable "admin_role_arn_mapping_template" {
  description = "The template that renders into yaml for the aws iam authenticator. Whitespace is important here."
  type        = "string"

  default = <<TEMPLATE
    - rolearn: %s
      username: admin
      groups:
      - system:masters
TEMPLATE
}

variable "sre_role_arns" {
  description = "A list of ARNs that will be mapped to cluster sre sre-administrators"
  type        = "list"
}

variable "sre_role_arn_mapping_template" {
  description = "The template that renders into yaml for the aws iam authenticator. Whitespace is important here."
  type        = "string"

  default = <<TEMPLATE
    - rolearn: %s
      username: sre
      groups:
      - sre
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
    - rolearn: %s
      username: dev
      groups:
      - dev
TEMPLATE
}

variable "bootstrapper_role_arn_mapping_template" {
  description = "The template that renders into yaml for the aws iam authenticator. Whitespace is important here."
  type        = "string"

  default = <<TEMPLATE
    - rolearn: %s
      username: system:node:{{EC2PrivateDNSName}}
      groups:
        - system:bootstrappers
        - system:nodes
TEMPLATE
}
