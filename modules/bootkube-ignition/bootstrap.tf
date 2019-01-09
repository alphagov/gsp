data "template_file" "bootstrap-apiserver" {
  template = "${file("${path.module}/data/bootstrap-manifests/bootstrap-apiserver.yaml")}"

  vars {
    etcd_servers = "${join(",", formatlist("https://%s:%s", var.etcd_servers, var.etcd_port))}"
    service_cidr = "${var.service_cidr}"
    k8s_tag      = "${var.k8s_tag}"
  }
}

data "ignition_file" "bootstrap-apiserver" {
  filesystem = "root"
  path       = "${var.assets_dir}/bootstrap-manifests/bootstrap-apiserver.yaml"
  mode       = 416

  content {
    content = "${data.template_file.bootstrap-apiserver.rendered}"
  }
}

data "template_file" "bootstrap-controller-manager" {
  template = "${file("${path.module}/data/bootstrap-manifests/bootstrap-controller-manager.yaml")}"

  vars {
    pod_cidr     = "${var.pod_cidr}"
    service_cidr = "${var.service_cidr}"
    k8s_tag      = "${var.k8s_tag}"
  }
}

data "ignition_file" "bootstrap-controller-manager" {
  filesystem = "root"
  path       = "${var.assets_dir}/bootstrap-manifests/bootstrap-controller-manager.yaml"
  mode       = 416

  content {
    content = "${data.template_file.bootstrap-controller-manager.rendered}"
  }
}

data "template_file" "bootstrap-scheduler" {
  template = "${file("${path.module}/data/bootstrap-manifests/bootstrap-scheduler.yaml")}"

  vars {
    k8s_tag = "${var.k8s_tag}"
  }
}

data "ignition_file" "bootstrap-scheduler" {
  filesystem = "root"
  path       = "${var.assets_dir}/bootstrap-manifests/bootstrap-scheduler.yaml"
  mode       = 416

  content {
    content = "${data.template_file.bootstrap-scheduler.rendered}"
  }
}
