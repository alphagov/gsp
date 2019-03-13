resource "tls_private_key" "sealed-secrets-key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "tls_self_signed_cert" "sealed-secrets-certificate" {
  key_algorithm   = "${tls_private_key.sealed-secrets-key.algorithm}"
  private_key_pem = "${tls_private_key.sealed-secrets-key.private_key_pem}"

  subject {
    common_name  = "${var.cluster_name}.${var.dns_zone}"
    organization = "Government Digital Service"
  }

  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
  ]
}
