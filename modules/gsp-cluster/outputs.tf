output "kubeconfig" {
  value = module.k8s-cluster.kubeconfig
}

output "worker_security_group_id" {
  value = module.k8s-cluster.worker_security_group_id
}

output "oidc_provider_url" {
  value = module.k8s-cluster.oidc_provider_url
}

output "oidc_provider_arn" {
  value = module.k8s-cluster.oidc_provider_arn
}

output "values" {
  sensitive = true
  value     = data.template_file.values.rendered
}
