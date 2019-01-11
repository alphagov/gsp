data "template_file" "flux" {
  template = "${file("${path.module}/data/flux.yaml")}"

  vars {
    namespace = "flux-system"
  }
}

resource "local_file" "flux" {
  filename = "${var.addons_dir}/flux.yaml"
  content  = "${data.template_file.flux.rendered}"
}

data "template_file" "namespace" {
  template = "${file("${path.module}/data/namespace.yaml")}"

  vars {
    namespace = "${var.namespace}"
  }
}

resource "local_file" "namespace" {
  filename = "${var.addons_dir}/${var.namespace}-namespace.yaml"
  content  = "${data.template_file.namespace.rendered}"
}

data "template_file" "helm-release" {
  template = "${file("${path.module}/data/helm-release.yaml")}"

  vars {
    namespace  = "${var.namespace}"
    chart_git  = "${var.chart_git}"
    chart_ref  = "${var.chart_ref}"
    chart_path = "${var.chart_path}"
    cluster_name = "${var.cluster_name}"
    cluster_domain = "${var.cluster_domain}"
    values     = "${var.values}"
    valueFileSecrets = "[${join(",",formatlist("{\"name\":\"%s\"}", var.valueFileSecrets))}]"
  }
}

resource "local_file" "helm-release-yaml" {
  filename = "${var.addons_dir}/${var.namespace}-helm-release.yaml"
  content  = "${data.template_file.helm-release.rendered}"
}

data "template_file" "values" {
  count = "${length(var.valueFileSecrets)}"
  template = "${file("${path.module}/data/values-secret.yaml")}"

  vars {
    namespace = "${var.namespace}"
    name = "${var.valueFileSecrets[count.index]}"
  }
}

resource "local_file" "values" {
  count = "${length(var.valueFileSecrets)}"
  filename = "${var.addons_dir}/${var.namespace}-${var.valueFileSecrets[count.index]}-values.yaml"
  content  = "${data.template_file.values.*.rendered[count.index]}"
}

output "release-name" {
  value = "${var.namespace}"
}
