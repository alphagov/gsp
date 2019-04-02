resource "random_string" "concourse_password" {
  length = 64
}

resource "random_string" "notary_passphrase_root" {
  length = 64
}

resource "random_string" "notary_passphrase_targets" {
  length = 64
}

resource "random_string" "notary_passphrase_snapshot" {
  length = 64
}

resource "random_string" "notary_passphrase_delegation" {
  length = 64
}

resource "random_string" "harbor_password" {
  length = 16
}

resource "random_string" "harbor_secret_key" {
  length = 16
}

resource "tls_private_key" "notary_root_key" {
  algorithm = "RSA"
  rsa_bits  = "4096"
}

resource "tls_private_key" "notary_ci_key" {
  algorithm = "RSA"
  rsa_bits  = "4096"
}
