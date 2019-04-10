resource "local_file" "tiller-role-binding" {
  filename = "addons/${var.cluster_name}/tiller-role-binding.yaml"
  content  = "${file("${path.module}/data/tiller-role-binding.yaml")}"
}

resource "local_file" "tiller-sa" {
  filename = "addons/${var.cluster_name}/tiller-sa.yaml"
  content  = "${file("${path.module}/data/tiller-sa.yaml")}"
}

resource "local_file" "tiller-svc" {
  filename = "addons/${var.cluster_name}/tiller-svc.yaml"
  content  = "${file("${path.module}/data/tiller-svc.yaml")}"
}

resource "local_file" "tiller" {
  filename = "addons/${var.cluster_name}/tiller.yaml"
  content  = "${file("${path.module}/data/tiller.yaml")}"
}
