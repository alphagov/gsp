data "aws_iam_policy_document" "kiam_server_role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "AWS"
      identifiers = ["${module.k8s-cluster.kiam-server-node-instance-role-arn}"]
    }
  }
}

data "aws_iam_policy_document" "kiam_server_policy" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    resources = [
      "${aws_iam_role.cloudwatch_log_shipping_role.arn}",
      "${aws_iam_role.concourse.arn}",
    ]
  }
}

# a trust relationship policy for roles we want kiam-server to assume
data "aws_iam_policy_document" "trust_kiam_server" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "AWS"
      identifiers = ["${aws_iam_role.kiam_server_role.arn}"]
    }
  }
}

resource "aws_iam_role" "kiam_server_role" {
  name        = "${var.cluster_name}_kiam_server"
  description = "Role the Kiam Server process assumes"

  assume_role_policy = "${data.aws_iam_policy_document.kiam_server_role.json}"
}

resource "aws_iam_policy" "kiam_server_policy" {
  name        = "${var.cluster_name}_kiam_server_policy"
  description = "Policy for the Kiam Server process"

  policy = "${data.aws_iam_policy_document.kiam_server_policy.json}"
}

resource "aws_iam_policy_attachment" "kiam_server_policy_attach" {
  name       = "${var.cluster_name}_kiam-server-attachment"
  roles      = ["${aws_iam_role.kiam_server_role.name}"]
  policy_arn = "${aws_iam_policy.kiam_server_policy.arn}"
}

resource "tls_private_key" "kiam_ca" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "kiam_ca" {
  key_algorithm   = "${tls_private_key.kiam_ca.algorithm}"
  private_key_pem = "${tls_private_key.kiam_ca.private_key_pem}"

  subject {
    common_name = "Kiam CA"
  }

  is_ca_certificate     = true
  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "cert_signing",
  ]
}

resource "tls_private_key" "kiam_server" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_cert_request" "kiam_server" {
  key_algorithm   = "${tls_private_key.kiam_server.algorithm}"
  private_key_pem = "${tls_private_key.kiam_server.private_key_pem}"

  subject {
    common_name = "gsp-kiam-server"
  }

  dns_names = [
    "127.0.0.1:443",
    "127.0.0.1:9610",
    "gsp-kiam-server",
    "gsp-kiam-server:443",
  ]

  ip_addresses = [
    "127.0.0.1",
  ]
}

resource "tls_locally_signed_cert" "kiam_server" {
  cert_request_pem      = "${tls_cert_request.kiam_server.cert_request_pem}"
  ca_key_algorithm      = "${tls_private_key.kiam_ca.algorithm}"
  ca_private_key_pem    = "${tls_private_key.kiam_ca.private_key_pem}"
  ca_cert_pem           = "${tls_self_signed_cert.kiam_ca.cert_pem}"
  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
    "client_auth",
  ]
}

resource "tls_private_key" "kiam_agent" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_cert_request" "kiam_agent" {
  key_algorithm   = "${tls_private_key.kiam_agent.algorithm}"
  private_key_pem = "${tls_private_key.kiam_agent.private_key_pem}"

  subject {
    common_name = "kiam agent"
  }
}

resource "tls_locally_signed_cert" "kiam_agent" {
  cert_request_pem      = "${tls_cert_request.kiam_agent.cert_request_pem}"
  ca_key_algorithm      = "${tls_private_key.kiam_ca.algorithm}"
  ca_private_key_pem    = "${tls_private_key.kiam_ca.private_key_pem}"
  ca_cert_pem           = "${tls_self_signed_cert.kiam_ca.cert_pem}"
  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
    "client_auth",
  ]
}
