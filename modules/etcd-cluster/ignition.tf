data "template_file" "etcds" {
  count    = "${var.node_count}"
  template = "etcd$${index}=https://$${cluster_name}-etcd$${index}.$${dns_zone}:2380"

  vars {
    index        = "${count.index}"
    cluster_name = "${var.cluster_name}"
    dns_zone     = "${var.dns_zone}"
  }
}

data "template_file" "user_data_object_key" {
  count    = "${var.node_count}"
  template = "/user_data/$${cluster_name}-etcd$${index}.userdata"

  vars {
    index        = "${count.index}"
    cluster_name = "${var.cluster_name}"
  }
}

data "template_file" "etcd_service" {
  count    = "${var.node_count}"
  template = "${file("${path.module}/data/etcd.service")}"

  vars {
    etcd_name            = "etcd${count.index}"
    etcd_domain          = "${var.cluster_name}-etcd${count.index}.${var.dns_zone}"
    etcd_initial_cluster = "${join(",", data.template_file.etcds.*.rendered)}"
  }
}

data "ignition_systemd_unit" "etcd_service" {
  count = "${var.node_count}"
  name  = "etcd-member.service"

  dropin {
    name    = "40-etcd-cluster.conf"
    content = "${element(data.template_file.etcd_service.*.rendered, count.index)}"
  }
}

data "ignition_systemd_unit" "wait-for-dns-service" {
  name = "wait-for-dns.service"

  content = "${file("${path.module}/data/wait-for-dns.service")}"
}

data "ignition_config" "etcd" {
  count = "${var.node_count}"

  files = [
    "${data.ignition_file.etcd-etcd-server-ca-crt.id}",
    "${data.ignition_file.etcd-etcd-server-key.id}",
    "${data.ignition_file.etcd-etcd-server-crt.id}",
    "${data.ignition_file.etcd-etcd-peer-ca-crt.id}",
    "${data.ignition_file.etcd-etcd-peer-key.id}",
    "${data.ignition_file.etcd-etcd-peer-crt.id}",
    "${data.ignition_file.etcd-etcd-client-ca-crt.id}",
    "${data.ignition_file.etcd-etcd-client-key.id}",
    "${data.ignition_file.etcd-etcd-client-crt.id}",
  ]

  systemd = ["${concat(
        list(
           element(data.ignition_systemd_unit.etcd_service.*.id, count.index),
           data.ignition_systemd_unit.wait-for-dns-service.id,
        ),
        module.common.ignition-systemd-unit-ids,
    )}"]
}

resource "aws_s3_bucket_object" "etcd-user-data" {
  count = "${var.node_count}"

  bucket  = "${var.user_data_bucket_name}"
  key     = "${element(data.template_file.user_data_object_key.*.rendered, count.index)}"
  content = "${element(data.ignition_config.etcd.*.rendered, count.index)}"
}

data "ignition_config" "etcd-actual" {
  count = "${var.node_count}"

  replace = {
    source       = "s3://${var.user_data_bucket_name}${element(data.template_file.user_data_object_key.*.rendered, count.index)}"
    verification = "sha512-${sha512(element(data.ignition_config.etcd.*.rendered, count.index))}"
  }
}
