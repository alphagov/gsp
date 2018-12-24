data "template_file" "kubelet-service" {
  template = "${file("${path.module}/data/kubelet.service")}"

  vars {
    dns_service_ip = "${var.dns_service_ip}"
    node_labels    = "${var.node_labels}"
    node_taints    = "${var.node_taints}"
    cluster_domain = "${var.cluster_domain_suffix}"
    k8s_tag        = "${var.k8s_tag}"
  }
}

data "ignition_file" "kubelet-kubeconfig" {
  filesystem = "root"
  path       = "/etc/kubernetes/kubeconfig"
  mode       = 416

  content {
    content = "${var.kubelet_kubeconfig}"
  }
}

data "ignition_file" "kube-ca-crt" {
  filesystem = "root"
  path       = "/etc/kubernetes/ca.crt"
  mode       = 416

  content {
    content = "${var.kube_ca_crt}"
  }
}

data "ignition_systemd_unit" "kubelet-service" {
  name    = "kubelet.service"
  content = "${data.template_file.kubelet-service.rendered}"
}

data "ignition_systemd_unit" "wait-for-dns-service" {
  name = "wait-for-dns.service"

  content = "${file("${path.module}/data/wait-for-dns.service")}"
}
