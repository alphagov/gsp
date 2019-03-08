data "template_file" "aws-iam-authenticator-cfg" {
  template = "${file("${path.module}/data/manifests/aws-iam-authenticator-cfg.yaml")}"

  vars {
    cluster_id              = "${var.cluster_id}"
    iam_admin_role_mappings = "${join("\n", formatlist(var.admin_role_arn_mapping_template, var.admin_role_arns))}"
    iam_sre_role_mappings   = "${join("\n", formatlist(var.sre_role_arn_mapping_template, var.sre_role_arns))}"
    iam_dev_role_mappings   = "${join("\n", formatlist(var.dev_role_arn_mapping_template, var.dev_role_arns))}"
  }
}

data "ignition_file" "aws-iam-authenticator-cfg" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/aws-iam-authenticator-cfg.yaml"
  mode       = 416

  content {
    content = "${data.template_file.aws-iam-authenticator-cfg.rendered}"
  }
}

data "ignition_file" "aws-iam-authenticator-daemonset" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/aws-iam-authenticator-daemonset.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/aws-iam-authenticator-daemonset.yaml")}"
  }
}

data "template_file" "aws-iam-authenticator-kubeconfig" {
  template = "${file("${path.module}/data/manifests/aws-iam-authenticator-kubeconfig.yaml")}"

  vars {
    iam_authenticator_cert = "${base64encode(tls_locally_signed_cert.iam-authenticator.cert_pem)}"
  }
}

data "ignition_file" "aws-iam-authenticator-kubeconfig" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/aws-iam-authenticator-kubeconfig.yaml"
  mode       = 416

  content {
    content = "${data.template_file.aws-iam-authenticator-kubeconfig.rendered}"
  }
}

data "template_file" "aws-iam-authenticator-secret" {
  template = "${file("${path.module}/data/manifests/aws-iam-authenticator-secret.yaml")}"

  vars {
    iam_authenticator_cert = "${base64encode(tls_locally_signed_cert.iam-authenticator.cert_pem)}"
    iam_authenticator_key  = "${base64encode(tls_private_key.iam-authenticator.private_key_pem)}"
  }
}

data "ignition_file" "aws-iam-authenticator-secret" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/aws-iam-authenticator-secret.yaml"
  mode       = 416

  content {
    content = "${data.template_file.aws-iam-authenticator-secret.rendered}"
  }
}

data "ignition_file" "default-storage-class" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/default-storage-class.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/default-storage-class.yaml")}"
  }
}

data "template_file" "kube-apiserver-secret" {
  template = "${file("${path.module}/data/manifests/kube-apiserver-secret.yaml")}"

  vars {
    apiserver_crt       = "${base64encode(tls_locally_signed_cert.apiserver.cert_pem)}"
    apiserver_key       = "${base64encode(tls_private_key.apiserver.private_key_pem)}"
    ca_data             = "${base64encode(tls_self_signed_cert.kube-ca.cert_pem)}"
    etcd_ca_crt         = "${base64encode(var.etcd_ca_cert_pem)}"
    etcd_client_crt     = "${base64encode(var.etcd_client_cert_pem)}"
    etcd_client_key     = "${base64encode(var.etcd_client_private_key_pem)}"
    service_account_pub = "${base64encode(tls_private_key.service-account.public_key_pem)}"
  }
}

data "ignition_file" "kube-apiserver-secret" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/kube-apiserver-secret.yaml"
  mode       = 416

  content {
    content = "${data.template_file.kube-apiserver-secret.rendered}"
  }
}

data "template_file" "kube-apiserver" {
  template = "${file("${path.module}/data/manifests/kube-apiserver.yaml")}"

  vars {
    etcd_servers      = "${join(",", formatlist("https://%s:%s", var.etcd_servers, var.etcd_port))}"
    service_cidr      = "${var.service_cidr}"
    apiserver_address = "${var.apiserver_address}"
    k8s_tag           = "${var.k8s_tag}"
  }
}

data "ignition_file" "kube-apiserver" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/kube-apiserver.yaml"
  mode       = 416

  content {
    content = "${data.template_file.kube-apiserver.rendered}"
  }
}

data "ignition_file" "kube-controller-manager-disruption" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/kube-controller-manager-disruption.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/kube-controller-manager-disruption.yaml")}"
  }
}

data "ignition_file" "kube-controller-manager-role-binding" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/kube-controller-manager-role-binding.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/kube-controller-manager-role-binding.yaml")}"
  }
}

data "template_file" "kube-controller-manager-secret" {
  template = "${file("${path.module}/data/manifests/kube-controller-manager-secret.yaml")}"

  vars {
    ca_crt              = "${base64encode(tls_self_signed_cert.kube-ca.cert_pem)}"
    ca_key              = "${base64encode(tls_private_key.kube-ca.private_key_pem)}"
    service_account_key = "${base64encode(tls_private_key.service-account.private_key_pem)}"
  }
}

data "ignition_file" "kube-controller-manager-secret" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/kube-controller-manager-secret.yaml"
  mode       = 416

  content {
    content = "${data.template_file.kube-controller-manager-secret.rendered}"
  }
}

data "ignition_file" "kube-controller-manager-service-account" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/kube-controller-manager-service-account.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/kube-controller-manager-service-account.yaml")}"
  }
}

data "template_file" "kube-controller-manager" {
  template = "${file("${path.module}/data/manifests/kube-controller-manager.yaml")}"

  vars {
    pod_cidr     = "${var.pod_cidr}"
    service_cidr = "${var.service_cidr}"
    k8s_tag      = "${var.k8s_tag}"
  }
}

data "ignition_file" "kube-controller-manager" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/kube-controller-manager.yaml"
  mode       = 416

  content {
    content = "${data.template_file.kube-controller-manager.rendered}"
  }
}

data "ignition_file" "kube-proxy-role-binding" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/kube-proxy-role-binding.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/kube-proxy-role-binding.yaml")}"
  }
}

data "ignition_file" "kube-proxy-sa" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/kube-proxy-sa.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/kube-proxy-sa.yaml")}"
  }
}

data "template_file" "kube-proxy" {
  template = "${file("${path.module}/data/manifests/kube-proxy.yaml")}"

  vars {
    k8s_tag  = "${var.k8s_tag}"
    pod_cidr = "${var.pod_cidr}"
  }
}

data "ignition_file" "kube-proxy" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/kube-proxy.yaml"
  mode       = 416

  content {
    content = "${data.template_file.kube-proxy.rendered}"
  }
}

data "ignition_file" "kube-scheduler-disruption" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/kube-scheduler-disruption.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/kube-scheduler-disruption.yaml")}"
  }
}

data "template_file" "kube-scheduler" {
  template = "${file("${path.module}/data/manifests/kube-scheduler.yaml")}"

  vars {
    k8s_tag = "${var.k8s_tag}"
  }
}

data "ignition_file" "kube-scheduler" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/kube-scheduler.yaml"
  mode       = 416

  content {
    content = "${data.template_file.kube-scheduler.rendered}"
  }
}

data "ignition_file" "kube-system-rbac-role-binding" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/kube-system-rbac-role-binding.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/kube-system-rbac-role-binding.yaml")}"
  }
}

data "template_file" "kubeconfig-in-cluster" {
  template = "${file("${path.module}/data/manifests/kubeconfig-in-cluster.yaml")}"

  vars {
    apiserver_address = "${var.apiserver_address}"
  }
}

data "ignition_file" "kubeconfig-in-cluster" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/kubeconfig-in-cluster.yaml"
  mode       = 416

  content {
    content = "${data.template_file.kubeconfig-in-cluster.rendered}"
  }
}

data "ignition_file" "kubernetes-dashboard-role-binding" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/kubernetes-dashboard-role-binding.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/kubernetes-dashboard-role-binding.yaml")}"
  }
}

data "ignition_file" "kubernetes-dashboard-role" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/kubernetes-dashboard-role.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/kubernetes-dashboard-role.yaml")}"
  }
}

data "ignition_file" "kubernetes-dashboard-sa" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/kubernetes-dashboard-sa.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/kubernetes-dashboard-sa.yaml")}"
  }
}

data "ignition_file" "kubernetes-dashboard-secret" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/kubernetes-dashboard-secret.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/kubernetes-dashboard-secret.yaml")}"
  }
}

data "ignition_file" "kubernetes-dashboard-svc" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/kubernetes-dashboard-svc.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/kubernetes-dashboard-svc.yaml")}"
  }
}

data "ignition_file" "kubernetes-dashboard" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/kubernetes-dashboard.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/kubernetes-dashboard.yaml")}"
  }
}

data "ignition_file" "pod-checkpointer-role-binding" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/pod-checkpointer-role-binding.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/pod-checkpointer-role-binding.yaml")}"
  }
}

data "ignition_file" "pod-checkpointer-role" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/pod-checkpointer-role.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/pod-checkpointer-role.yaml")}"
  }
}

data "ignition_file" "pod-checkpointer-sa" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/pod-checkpointer-sa.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/pod-checkpointer-sa.yaml")}"
  }
}

data "ignition_file" "pod-checkpointer" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/pod-checkpointer.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/pod-checkpointer.yaml")}"
  }
}

data "ignition_file" "tiller-role-binding" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/tiller-role-binding.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/tiller-role-binding.yaml")}"
  }
}

data "ignition_file" "tiller-sa" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/tiller-sa.yaml"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/tiller-sa.yaml")}"
  }
}

data "ignition_file" "tiller-svc" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/tiller-svc"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/tiller-svc.yaml")}"
  }
}

data "ignition_file" "tiller" {
  filesystem = "root"
  path       = "${var.assets_dir}/manifests/tiller"
  mode       = 416

  content {
    content = "${file("${path.module}/data/manifests/tiller.yaml")}"
  }
}
