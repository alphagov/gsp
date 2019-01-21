data "template_file" "role" {
  count    = "${length(var.namespaces)}"
  template = "${file("${path.module}/data/role.yaml")}"

  vars = {
    namespace = "${element(var.namespaces, count.index)}"
  }
}

resource "local_file" "role" {
  count    = "${length(var.namespaces)}"
  filename = "${var.addons_dir}/${element(var.namespaces, count.index)}-role.yaml"
  content  = "${element(data.template_file.role.*.rendered, count.index)}"
}

data "template_file" "role-binding" {
  count    = "${length(var.namespaces)}"
  template = "${file("${path.module}/data/role-binding.yaml")}"

  vars = {
    namespace  = "${element(var.namespaces, count.index)}"
    group_name = "${var.group_name}"
  }
}

resource "local_file" "role-binding" {
  count    = "${length(var.namespaces)}"
  filename = "${var.addons_dir}/${element(var.namespaces, count.index)}-role-binding.yaml"
  content  = "${element(data.template_file.role-binding.*.rendered, count.index)}"
}
