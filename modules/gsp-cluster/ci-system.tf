data "aws_iam_policy_document" "assume-harbor" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals = {
      type        = "AWS"
      identifiers = ["${aws_iam_role.kiam_server_role.arn}"]
    }
  }
}

resource "aws_iam_role" "harbor" {
  name               = "${var.cluster_name}-harbor"
  description        = "Role the harbor process assumes"
  assume_role_policy = "${data.aws_iam_policy_document.assume-harbor.json}"
}

data "aws_iam_policy_document" "harbor-s3" {
  statement {
    actions = [
      "s3:*",
    ]

    resources = [
      "${element(concat(aws_s3_bucket.ci-system-harbor-registry-storage.*.arn, list("")), 0)}",
      "${element(concat(aws_s3_bucket.ci-system-harbor-registry-storage.*.arn, list("")), 0)}/*",
    ]
  }
}

resource "aws_iam_policy" "harbor-s3" {
  name        = "${var.cluster_name}-harbor-s3"
  description = "Policy for the harbor s3 access"
  policy      = "${data.aws_iam_policy_document.harbor-s3.json}"
}

resource "aws_iam_policy_attachment" "harbor-s3" {
  name       = "${var.cluster_name}-harbor-s3"
  roles      = ["${element(aws_iam_role.harbor.*.name, count.index)}"]
  policy_arn = "${element(aws_iam_policy.harbor-s3.*.arn, count.index)}"
}

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

resource "aws_s3_bucket" "ci-system-harbor-registry-storage" {
  bucket = "registry-${var.cluster_name}-${var.account_name}"
  acl    = "private"

  force_destroy = true # NEED TO VALIDATE!!!

  tags = {
    Name = "Harbor registry and chartmuseum storage"
  }
}
