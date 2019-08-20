data "template_file" "istio_values" {
  template = "${file("${path.module}/data/gsp-istio-values.yaml")}"

  vars {
    cert_manager_role_name = "${aws_iam_role.cert_manager.name}"
  }
}
