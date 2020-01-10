resource "aws_iam_role" "harbor" {
  name               = "${var.cluster_name}-harbor"
  description        = "Role the harbor process assumes"
  assume_role_policy = data.aws_iam_policy_document.trust_kiam_server.json
}

data "aws_iam_policy_document" "harbor-s3" {
  statement {
    actions = [
      "s3:*",
    ]

    resources = [
      aws_s3_bucket.ci-system-harbor-registry-storage.arn,
      "${aws_s3_bucket.ci-system-harbor-registry-storage.arn}/*",
    ]
  }
}

resource "aws_iam_policy" "harbor-s3" {
  name        = "${var.cluster_name}-harbor-s3"
  description = "Policy for the harbor s3 access"
  policy      = data.aws_iam_policy_document.harbor-s3.json
}

resource "aws_iam_policy_attachment" "harbor-s3" {
  name       = "${var.cluster_name}-harbor-s3"
  roles      = [aws_iam_role.harbor.name]
  policy_arn = aws_iam_policy.harbor-s3.arn
}

resource "random_password" "concourse_password" {
  length  = 64
  special = false
}

resource "random_password" "notary_passphrase_root" {
  length  = 64
  special = false
}

resource "random_password" "notary_passphrase_targets" {
  length  = 64
  special = false
}

resource "random_password" "notary_passphrase_snapshot" {
  length  = 64
  special = false
}

resource "random_password" "notary_passphrase_delegation" {
  length  = 64
  special = false
}

resource "random_password" "harbor_password" {
  length  = 16
  special = false
}

resource "random_password" "harbor_secret_key" {
  length  = 16
  special = false
}

resource "tls_private_key" "notary_root_key" {
  algorithm = "RSA"
  rsa_bits  = "4096"
}

resource "tls_private_key" "notary_ci_key" {
  algorithm = "RSA"
  rsa_bits  = "4096"
}

resource "tls_self_signed_cert" "notary_root_ca" {
  key_algorithm   = tls_private_key.notary_root_key.algorithm
  private_key_pem = tls_private_key.notary_root_key.private_key_pem

  subject {
    common_name  = "gsp-harbor-notary-signer"
    organization = "gsp"
  }

  is_ca_certificate     = true
  validity_period_hours = 26280 # 3yrs

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
    "cert_signing",
  ]
}

resource "tls_cert_request" "notary_cert" {
  key_algorithm   = tls_private_key.notary_root_key.algorithm
  private_key_pem = tls_private_key.notary_root_key.private_key_pem

  subject {
    common_name  = "gsp-harbor-notary-signer"
    organization = "gsp"
  }
}

resource "tls_locally_signed_cert" "notary_cert" {
  cert_request_pem   = tls_cert_request.notary_cert.cert_request_pem
  ca_key_algorithm   = tls_private_key.notary_root_key.algorithm
  ca_private_key_pem = tls_private_key.notary_root_key.private_key_pem
  ca_cert_pem        = tls_self_signed_cert.notary_root_ca.cert_pem

  validity_period_hours = 8760 # 1yr

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
    "cert_signing",
  ]
}

resource "aws_s3_bucket" "ci-system-harbor-registry-storage" {
  bucket = "registry-${var.cluster_name}-${var.account_name}"
  acl    = "private"

  force_destroy = true # NEED TO VALIDATE!!!

  tags = {
    Name = "Harbor registry and chartmuseum storage"
  }
}

resource "random_password" "harbor_db_master_password" {
  length  = 32
  special = false
}

resource "aws_rds_cluster" "harbor" {
  cluster_identifier = "${var.cluster_name}-harbor"
  master_password    = random_password.harbor_db_master_password.result
  master_username    = "harbor"
  availability_zones = var.availability_zones
  engine             = "aurora-postgresql"
  engine_version     = "10.7"
  enabled_cloudwatch_logs_exports = [
    "postgresql",
  ]
  db_subnet_group_name      = aws_db_subnet_group.private.name
  vpc_security_group_ids    = [aws_security_group.rds-from-worker.id]
  final_snapshot_identifier = "${var.cluster_name}-harbor-final"
  skip_final_snapshot       = "${var.harbor_rds_skip_final_snapshot}"
  tags = {
    Name = var.cluster_name
  }
}

resource "aws_rds_cluster_instance" "harbor" {
  count                = 2
  identifier           = "${var.cluster_name}-harbor-${count.index}"
  cluster_identifier   = aws_rds_cluster.harbor.id
  engine               = "aurora-postgresql"
  engine_version       = "10.7"
  apply_immediately    = true
  ca_cert_identifier   = "rds-ca-2019"
  instance_class       = "db.t3.medium"
  db_subnet_group_name = aws_db_subnet_group.private.name

  tags = {
    Name = var.cluster_name
  }
}
