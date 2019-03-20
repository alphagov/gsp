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

resource "aws_iam_role" "sre" {
  name = "${var.cluster_name}-sre"

  assume_role_policy = "${data.aws_iam_policy_document.grant-iam-sre-policy.json}"
}

data "aws_iam_policy_document" "grant-iam-sre-policy" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals = {
      type = "AWS"

      identifiers = ["${concat(var.admin_role_arns, var.sre_user_arns)}"]
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

data "aws_iam_policy_document" "grant-iam-dev" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals = {
      type        = "AWS"
      identifiers = ["${concat(var.admin_role_arns, var.dev_user_arns)}"]
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

data "aws_iam_policy_document" "kiam_server_role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals = {
      type        = "AWS"
      identifiers = ["${module.k8s-cluster.controller-role-arn}"]
    }
  }
}

data "aws_iam_policy_document" "kiam_server_policy" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    resources = [
      "${aws_iam_role.cloudwatch_log_shipping_role.arn}",
      "${module.gsp-canary.canary_role_arn}",
    ]
  }
}

data "aws_iam_policy_document" "cloudwatch_log_shipping_role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals = {
      type        = "AWS"
      identifiers = ["${aws_iam_role.kiam_server_role.arn}"]
    }
  }
}

data "aws_iam_policy_document" "cloudwatch_log_shipping_policy" {
  statement {
    effect = "Allow"

    actions = [
      "logs:DescribeLogGroups",
    ]

    resources = ["*"]
  }

  statement {
    effect = "Allow"

    actions = [
      "logs:DescribeLogStreams",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = ["${aws_cloudwatch_log_group.logs.arn}"]
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

resource "aws_iam_role" "cloudwatch_log_shipping_role" {
  name = "${var.cluster_name}_cloudwatch_log_shipping_role"

  assume_role_policy = "${data.aws_iam_policy_document.cloudwatch_log_shipping_role.json}"
}

resource "aws_iam_policy" "cloudwatch_log_shipping_policy" {
  name        = "${var.cluster_name}_cloudwatch_log_shipping_policy"
  description = "Send logs to Clouwatch"

  policy = "${data.aws_iam_policy_document.cloudwatch_log_shipping_policy.json}"
}

resource "aws_iam_policy_attachment" "cloudwatch_log_shipping_policy" {
  name       = "${var.cluster_name}_cloudwatch_log_shipping_role_policy_attachement"
  roles      = ["${aws_iam_role.cloudwatch_log_shipping_role.name}"]
  policy_arn = "${aws_iam_policy.cloudwatch_log_shipping_policy.arn}"
}
