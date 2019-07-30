variable "namespace_name" {
  description = "unique name for the namespace"
}

variable "cluster_name" {
  description = "cluster name to scope this role to"
}

variable "role_name" {
  description = "name of the role"
}

variable "account_id" {
  description = "account ID this cluster is running in"
}