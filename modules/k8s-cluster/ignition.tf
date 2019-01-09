module "k8s-controller-node-ignition" {
  source                = "../k8s-node-ignition"
  dns_service_ip        = "${cidrhost(var.service_cidr, 10)}"
  node_labels           = "${var.controller_node_labels}"
  node_taints           = "${var.controller_node_taints}"
  cluster_domain_suffix = "${var.cluster_domain_suffix}"
  k8s_tag               = "${var.k8s_tag}"
  kubelet_kubeconfig    = "${var.kubelet_kubeconfig}"
  kube_ca_crt           = "${var.kube_ca_crt}"
}

module "k8s-worker-node-ignition" {
  source                = "../k8s-node-ignition"
  dns_service_ip        = "${cidrhost(var.service_cidr, 10)}"
  node_labels           = "${var.worker_node_labels}"
  node_taints           = "${var.worker_node_taints}"
  cluster_domain_suffix = "${var.cluster_domain_suffix}"
  k8s_tag               = "${var.k8s_tag}"
  kubelet_kubeconfig    = "${var.kubelet_kubeconfig}"
  kube_ca_crt           = "${var.kube_ca_crt}"
}

data "template_file" "controller-user-data-object-key" {
  template = "/user_data/$${cluster_name}-controller.userdata"

  vars {
    cluster_name = "${var.cluster_name}"
  }
}

data "template_file" "worker-user-data-object-key" {
  template = "/user_data/$${cluster_name}-worker.userdata"

  vars {
    cluster_name = "${var.cluster_name}"
  }
}

data "ignition_config" "controller" {
  files = ["${module.k8s-controller-node-ignition.ignition_file_ids}"]

  systemd = ["${concat(
        module.k8s-controller-node-ignition.ignition_systemd_unit_ids,
        module.common.ignition-systemd-unit-ids,
    )}"]
}

data "ignition_config" "worker" {
  files = ["${module.k8s-worker-node-ignition.ignition_file_ids}"]

  systemd = ["${concat(
        module.k8s-worker-node-ignition.ignition_systemd_unit_ids,
        module.common.ignition-systemd-unit-ids,
    )}"]
}

resource "aws_s3_bucket_object" "controller-user-data" {
  bucket  = "${var.user_data_bucket_name}"
  key     = "${data.template_file.controller-user-data-object-key.rendered}"
  content = "${data.ignition_config.controller.rendered}"
}

data "ignition_config" "controller-actual" {
  replace = {
    source       = "s3://${var.user_data_bucket_name}${aws_s3_bucket_object.controller-user-data.id}"
    verification = "sha512-${sha512(data.ignition_config.controller.rendered)}"
  }
}

resource "aws_s3_bucket_object" "worker-user-data" {
  bucket  = "${var.user_data_bucket_name}"
  key     = "${data.template_file.worker-user-data-object-key.rendered}"
  content = "${data.ignition_config.worker.rendered}"
}

data "ignition_config" "worker-actual" {
  replace = {
    source       = "s3://${var.user_data_bucket_name}${aws_s3_bucket_object.worker-user-data.id}"
    verification = "sha512-${sha512(data.ignition_config.worker.rendered)}"
  }
}
