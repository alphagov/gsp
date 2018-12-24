data "ignition_file" "calico-bgpconfigurations-crd" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/calico-bgpconfigurations-crd.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests-networking/calico-bgpconfigurations-crd.yaml")}"
  }
}

data "ignition_file" "calico-bgppeers-crd" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/calico-bgppeers-crd.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests-networking/calico-bgppeers-crd.yaml")}"
  }
}

data "ignition_file" "calico-clusterinformations-crd" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/calico-clusterinformations-crd.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests-networking/calico-clusterinformations-crd.yaml")}"
  }
}

data "ignition_file" "calico-config" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/calico-config.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests-networking/calico-config.yaml")}"
  }
}

data "ignition_file" "calico-felixconfigurations-crd" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/calico-felixconfigurations-crd.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests-networking/calico-felixconfigurations-crd.yaml")}"
  }
}

data "ignition_file" "calico-globalnetworkpolicies-crd" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/calico-globalnetworkpolicies-crd.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests-networking/calico-globalnetworkpolicies-crd.yaml")}"
  }
}

data "ignition_file" "calico-globalnetworksets-crd" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/calico-globalnetworksets-crd.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests-networking/calico-globalnetworksets-crd.yaml")}"
  }
}

data "ignition_file" "calico-hostendpoints-crd" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/calico-hostendpoints-crd.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests-networking/calico-hostendpoints-crd.yaml")}"
  }
}

data "ignition_file" "calico-ippools-crd" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/calico-ippools-crd.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests-networking/calico-ippools-crd.yaml")}"
  }
}

data "ignition_file" "calico-networkpolicies-crd" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/calico-networkpolicies-crd.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests-networking/calico-networkpolicies-crd.yaml")}"
  }
}

data "ignition_file" "calico-node-cluster-role-binding" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/calico-node-cluster-role-binding.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests-networking/calico-node-cluster-role-binding.yaml")}"
  }
}

data "ignition_file" "calico-node-cluster-role" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/calico-node-cluster-role.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests-networking/calico-node-cluster-role.yaml")}"
  }
}

data "ignition_file" "calico-node-service-account" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/calico-node-service-account.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests-networking/calico-node-service-account.yaml")}"
  }
}

data "template_file" "calico" {
  template = "${file("${path.module}/data/manifests-networking/calico.yaml")}"

  vars {
    pod_cidr = "${var.pod_cidr}"
  }
}

data "ignition_file" "calico" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/calico.yaml"
  mode       = 416

  content {
    content = "${data.template_file.calico.rendered}"
  }
}
