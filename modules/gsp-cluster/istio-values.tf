data "template_file" "istio-values" {
  template = "${file("${path.module}/data/istio-values.yaml")}"

  vars {
    cert_manager_role_name = "${aws_iam_role.cert-manager.name}"
  }
}
