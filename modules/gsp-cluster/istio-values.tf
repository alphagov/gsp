data "template_file" "istio-values" {
  template = "${file("${path.module}/data/istio-values.yaml")}"

  vars {}
}
