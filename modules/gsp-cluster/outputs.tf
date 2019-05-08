output "kubeconfig" {
  value = "${module.k8s-cluster.kubeconfig}"
}

output "istio-values" {
  sensitive = true
  value     = "${data.template_file.istio-values.rendered}"
}

output "values" {
  sensitive = true
  value     = "${data.template_file.values.rendered}"
}
