output "kubeconfig" {
  value = "${module.k8s-cluster.kubeconfig}"
}

output "values" {
  sensitive = true
  value     = "${data.template_file.values.rendered}"
}
