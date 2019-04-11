variable "name" {
  description = "name of namespace to create"
  type        = "string"
}

variable "cluster_name" {
  description = "name of this cluster/environment (accessible as .Values.cluster.name in charts)"
  type        = "string"
}

variable "cluster_domain" {
  description = "domain mapped to this cluster/environment (accessible as .Values.cluster.name in charts)"
  type        = "string"
}
