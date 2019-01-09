data "template_file" "bootkube-service" {
  template = "${file("${path.module}/data/bootkube.service")}"

  vars {
    assets_dir = "${var.assets_dir}"
  }
}

data "ignition_systemd_unit" "bootkube-service" {
  name    = "bootkube.service"
  content = "${data.template_file.bootkube-service.rendered}"
}
