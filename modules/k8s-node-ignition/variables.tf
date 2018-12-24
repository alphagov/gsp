variable "dns_service_ip" {
  type = "string"
}

variable "node_labels" {
  type = "string"
}

variable "node_taints" {
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
