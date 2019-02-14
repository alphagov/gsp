resource "tls_private_key" "etcd-ca" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "tls_self_signed_cert" "etcd-ca" {
  key_algorithm   = "${tls_private_key.etcd-ca.algorithm}"
  private_key_pem = "${tls_private_key.etcd-ca.private_key_pem}"

  subject {
    common_name  = "etcd-ca"
    organization = "etcd"
  }

  is_ca_certificate     = true
  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "cert_signing",
  ]
}

data "ignition_file" "etcd-client-ca-crt" {
  filesystem = "root"
  path       = "${var.assets_dir}/tls/etcd-client-ca.crt"
  mode       = 416

  content {
    content = "${tls_self_signed_cert.etcd-ca.cert_pem}"
  }
}

data "ignition_file" "etcd-etcd-client-ca-crt" {
  filesystem = "root"
  path       = "/etc/ssl/certs/etcd/etcd-client-ca.crt"
  mode       = 416

  content {
    content = "${tls_self_signed_cert.etcd-ca.cert_pem}"
  }
}

data "ignition_file" "etcd-server-ca-crt" {
  filesystem = "root"
  path       = "${var.assets_dir}/tls/etcd/server-ca.crt"
  mode       = 416

  content {
    content = "${tls_self_signed_cert.etcd-ca.cert_pem}"
  }
}

data "ignition_file" "etcd-etcd-server-ca-crt" {
  filesystem = "root"
  path       = "/etc/ssl/certs/etcd/server-ca.crt"
  mode       = 416

  content {
    content = "${tls_self_signed_cert.etcd-ca.cert_pem}"
  }
}

data "ignition_file" "etcd-peer-ca-crt" {
  filesystem = "root"
  path       = "${var.assets_dir}/tls/etcd/peer-ca.crt"
  mode       = 416

  content {
    content = "${tls_self_signed_cert.etcd-ca.cert_pem}"
  }
}

data "ignition_file" "etcd-etcd-peer-ca-crt" {
  filesystem = "root"
  path       = "/etc/ssl/certs/etcd/peer-ca.crt"
  mode       = 416

  content {
    content = "${tls_self_signed_cert.etcd-ca.cert_pem}"
  }
}

resource "tls_private_key" "etcd-client" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

data "ignition_file" "etcd-client-key" {
  filesystem = "root"
  path       = "${var.assets_dir}/tls/etcd-client.key"
  mode       = 384

  content {
    content = "${tls_private_key.etcd-client.private_key_pem}"
  }
}

data "ignition_file" "etcd-etcd-client-key" {
  filesystem = "root"
  path       = "/etc/ssl/certs/etcd/etcd-client.key"
  mode       = 384

  content {
    content = "${tls_private_key.etcd-client.private_key_pem}"
  }
}

resource "tls_cert_request" "etcd-client" {
  key_algorithm   = "${tls_private_key.etcd-client.algorithm}"
  private_key_pem = "${tls_private_key.etcd-client.private_key_pem}"

  subject {
    common_name  = "etcd-client"
    organization = "etcd"
  }

  ip_addresses = ["127.0.0.1"]

  dns_names = ["${concat(
    data.template_file.etcd_servers.*.rendered,
    list(
      "localhost",
    ))}"]
}

resource "tls_locally_signed_cert" "etcd-client" {
  cert_request_pem = "${tls_cert_request.etcd-client.cert_request_pem}"

  ca_key_algorithm   = "${tls_self_signed_cert.etcd-ca.key_algorithm}"
  ca_private_key_pem = "${tls_private_key.etcd-ca.private_key_pem}"
  ca_cert_pem        = "${tls_self_signed_cert.etcd-ca.cert_pem}"

  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
    "client_auth",
  ]
}

data "ignition_file" "etcd-client-crt" {
  filesystem = "root"
  path       = "${var.assets_dir}/tls/etcd-client.crt"
  mode       = 416

  content {
    content = "${tls_locally_signed_cert.etcd-client.cert_pem}"
  }
}

data "ignition_file" "etcd-etcd-client-crt" {
  filesystem = "root"
  path       = "/etc/ssl/certs/etcd/etcd-client.crt"
  mode       = 416

  content {
    content = "${tls_locally_signed_cert.etcd-client.cert_pem}"
  }
}

resource "tls_private_key" "etcd-server" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

data "ignition_file" "etcd-server-key" {
  filesystem = "root"
  path       = "${var.assets_dir}/tls/etcd/server.key"
  mode       = 384

  content {
    content = "${tls_private_key.etcd-server.private_key_pem}"
  }
}

data "ignition_file" "etcd-etcd-server-key" {
  filesystem = "root"
  path       = "/etc/ssl/certs/etcd/server.key"
  mode       = 384

  content {
    content = "${tls_private_key.etcd-server.private_key_pem}"
  }
}

resource "tls_cert_request" "etcd-server" {
  key_algorithm   = "${tls_private_key.etcd-server.algorithm}"
  private_key_pem = "${tls_private_key.etcd-server.private_key_pem}"

  subject {
    common_name  = "etcd-server"
    organization = "etcd"
  }

  ip_addresses = ["127.0.0.1"]

  dns_names = ["${concat(
    data.template_file.etcd_servers.*.rendered,
    list(
      "localhost",
    ))}"]
}

resource "tls_locally_signed_cert" "etcd-server" {
  cert_request_pem = "${tls_cert_request.etcd-server.cert_request_pem}"

  ca_key_algorithm   = "${tls_self_signed_cert.etcd-ca.key_algorithm}"
  ca_private_key_pem = "${tls_private_key.etcd-ca.private_key_pem}"
  ca_cert_pem        = "${tls_self_signed_cert.etcd-ca.cert_pem}"

  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
    "client_auth",
  ]
}

data "ignition_file" "etcd-server-crt" {
  filesystem = "root"
  path       = "${var.assets_dir}/tls/etcd/server.crt"
  mode       = 416

  content {
    content = "${tls_locally_signed_cert.etcd-server.cert_pem}"
  }
}

data "ignition_file" "etcd-etcd-server-crt" {
  filesystem = "root"
  path       = "/etc/ssl/certs/etcd/server.crt"
  mode       = 416

  content {
    content = "${tls_locally_signed_cert.etcd-server.cert_pem}"
  }
}

resource "tls_private_key" "etcd-peer" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

data "ignition_file" "etcd-peer-key" {
  filesystem = "root"
  path       = "${var.assets_dir}/tls/etcd/peer.key"
  mode       = 384

  content {
    content = "${tls_private_key.etcd-peer.private_key_pem}"
  }
}

data "ignition_file" "etcd-etcd-peer-key" {
  filesystem = "root"
  path       = "/etc/ssl/certs/etcd/peer.key"
  mode       = 384

  content {
    content = "${tls_private_key.etcd-peer.private_key_pem}"
  }
}

resource "tls_cert_request" "etcd-peer" {
  key_algorithm   = "${tls_private_key.etcd-peer.algorithm}"
  private_key_pem = "${tls_private_key.etcd-peer.private_key_pem}"

  subject {
    common_name  = "etcd-peer"
    organization = "etcd"
  }

  dns_names = ["${data.template_file.etcd_servers.*.rendered}"]
}

resource "tls_locally_signed_cert" "etcd-peer" {
  cert_request_pem = "${tls_cert_request.etcd-peer.cert_request_pem}"

  ca_key_algorithm   = "${tls_self_signed_cert.etcd-ca.key_algorithm}"
  ca_private_key_pem = "${tls_private_key.etcd-ca.private_key_pem}"
  ca_cert_pem        = "${tls_self_signed_cert.etcd-ca.cert_pem}"

  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
    "client_auth",
  ]
}

data "ignition_file" "etcd-peer-crt" {
  filesystem = "root"
  path       = "${var.assets_dir}/tls/etcd/peer.crt"
  mode       = 416

  content {
    content = "${tls_locally_signed_cert.etcd-peer.cert_pem}"
  }
}

data "ignition_file" "etcd-etcd-peer-crt" {
  filesystem = "root"
  path       = "/etc/ssl/certs/etcd/peer.crt"
  mode       = 416

  content {
    content = "${tls_locally_signed_cert.etcd-peer.cert_pem}"
  }
}
