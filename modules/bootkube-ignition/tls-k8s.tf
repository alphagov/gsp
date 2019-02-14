resource "tls_private_key" "kube-ca" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

data "ignition_file" "ca-key" {
  filesystem = "root"
  path       = "${var.assets_dir}/tls/ca.key"
  mode       = 384

  content {
    content = "${tls_private_key.kube-ca.private_key_pem}"
  }
}

resource "tls_self_signed_cert" "kube-ca" {
  key_algorithm   = "${tls_private_key.kube-ca.algorithm}"
  private_key_pem = "${tls_private_key.kube-ca.private_key_pem}"

  subject {
    common_name  = "kube-ca"
    organization = "bootkube"
  }

  is_ca_certificate     = true
  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "cert_signing",
  ]
}

data "ignition_file" "ca-crt" {
  filesystem = "root"
  path       = "${var.assets_dir}/tls/ca.crt"
  mode       = 416

  content {
    content = "${tls_self_signed_cert.kube-ca.cert_pem}"
  }
}

resource "tls_private_key" "apiserver" {
  algorithm = "${tls_private_key.kube-ca.algorithm}"
  rsa_bits  = 4096
}

data "ignition_file" "apiserver-key" {
  filesystem = "root"
  path       = "${var.assets_dir}/tls/apiserver.key"
  mode       = 384

  content {
    content = "${tls_private_key.apiserver.private_key_pem}"
  }
}

resource "tls_cert_request" "apiserver" {
  key_algorithm   = "${tls_private_key.kube-ca.algorithm}"
  private_key_pem = "${tls_private_key.apiserver.private_key_pem}"

  subject {
    common_name  = "kube-apiserver"
    organization = "system:masters"
  }

  dns_names = [
    "kubernetes",
    "kubernetes.default",
    "kubernetes.default.svc",
    "kubernetes.default.svc.${var.cluster_domain_suffix}",
    "${var.apiserver_address}",
  ]

  ip_addresses = ["${cidrhost(var.service_cidr, 1)}"]
}

resource "tls_locally_signed_cert" "apiserver" {
  cert_request_pem = "${tls_cert_request.apiserver.cert_request_pem}"

  ca_key_algorithm   = "${tls_self_signed_cert.kube-ca.key_algorithm}"
  ca_private_key_pem = "${tls_private_key.kube-ca.private_key_pem}"
  ca_cert_pem        = "${tls_self_signed_cert.kube-ca.cert_pem}"

  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
    "client_auth",
  ]
}

data "ignition_file" "apiserver-crt" {
  filesystem = "root"
  path       = "${var.assets_dir}/tls/apiserver.crt"
  mode       = 416

  content {
    content = "${tls_locally_signed_cert.apiserver.cert_pem}"
  }
}

resource "tls_private_key" "service-account" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

data "ignition_file" "service-account-key" {
  filesystem = "root"
  path       = "${var.assets_dir}/tls/service-account.key"
  mode       = 384

  content {
    content = "${tls_private_key.service-account.private_key_pem}"
  }
}

data "ignition_file" "service-account-pub" {
  filesystem = "root"
  path       = "${var.assets_dir}/tls/service-account.pub"
  mode       = 416

  content {
    content = "${tls_private_key.service-account.public_key_pem}"
  }
}

resource "tls_private_key" "kubelet" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

data "ignition_file" "admin-key" {
  filesystem = "root"
  path       = "${var.assets_dir}/tls/admin.key"
  mode       = 384

  content {
    content = "${tls_private_key.kubelet.private_key_pem}"
  }
}

resource "tls_cert_request" "kubelet" {
  key_algorithm   = "${tls_private_key.kubelet.algorithm}"
  private_key_pem = "${tls_private_key.kubelet.private_key_pem}"

  subject {
    common_name  = "kubelet"
    organization = "system:masters"
  }
}

resource "tls_locally_signed_cert" "kubelet" {
  cert_request_pem = "${tls_cert_request.kubelet.cert_request_pem}"

  ca_key_algorithm   = "${tls_self_signed_cert.kube-ca.key_algorithm}"
  ca_private_key_pem = "${tls_private_key.kube-ca.private_key_pem}"
  ca_cert_pem        = "${tls_self_signed_cert.kube-ca.cert_pem}"

  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
    "client_auth",
  ]
}

data "ignition_file" "admin-crt" {
  filesystem = "root"
  path       = "${var.assets_dir}/tls/admin.crt"
  mode       = 416

  content {
    content = "${tls_private_key.kubelet.public_key_pem}"
  }
}
