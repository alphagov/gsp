data "ignition_file" "coredns-cluster-role-binding" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/coredns-cluster-role-binding.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/coredns/cluster-role-binding.yaml")}"
  }
}

data "ignition_file" "coredns-cluster-role" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/coredns-cluster-role.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/coredns/cluster-role.yaml")}"
  }
}

data "template_file" "coredns-config-yaml" {
  template = "${file("${path.module}/data/manifests/coredns/config.yaml")}"

  vars {
    service_cidr   = "${var.service_cidr}"
    cluster_domain = "${var.cluster_domain_suffix}"
  }
}

data "ignition_file" "coredns-config-yaml" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/coredns-config.yaml"
  mode       = 416

  content {
    content = "${data.template_file.coredns-config-yaml.rendered}"
  }
}

data "ignition_file" "coredns-deployment" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/coredns-deployment.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/coredns/deployment.yaml")}"
  }
}

data "ignition_file" "coredns-service-account" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/coredns-service-account.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/coredns/service-account.yaml")}"
  }
}

data "template_file" "coredns-service" {
  template = "${file("${path.module}/data/manifests/coredns/service.yaml")}"

  vars {
    dns_service_ip = "${cidrhost(var.service_cidr, 10)}"
  }
}

data "ignition_file" "coredns-service" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/coredns-service.yaml"
  mode       = 416

  content {
    content = "${data.template_file.coredns-service.rendered}"
  }
}
