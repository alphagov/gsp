data "aws_iam_policy_document" "user_data_policy_document" {
  statement {
    actions = [
      "s3:GetObject",
    ]

    resources = [
      "${data.aws_s3_bucket.user_data.arn}/user_data/${var.cluster_name}-*",
    ]
  }
}

data "aws_s3_bucket" "user_data" {
  bucket = "${var.user_data_bucket_name}"
}

resource "aws_iam_policy" "s3-user-data-policy" {
  name   = "${var.cluster_name}-s3-user-data-policy"
  policy = "${data.aws_iam_policy_document.user_data_policy_document.json}"
}

resource "aws_iam_role" "dev" {
  name = "${var.cluster_name}-dev"

  assume_role_policy = "${data.aws_iam_policy_document.grant-iam-dev.json}"
}

data "aws_iam_policy_document" "grant-iam-dev" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals = {
      type        = "AWS"
      identifiers = ["${var.dev_user_arns}"]
    }

    condition {
      test     = "Bool"
      variable = "aws:MultiFactorAuthPresent"
      values   = ["true"]
    }

    condition {
      test     = "IpAddress"
      variable = "aws:SourceIp"
      values   = ["${var.gds_external_cidrs}"]
    }
  }
}
