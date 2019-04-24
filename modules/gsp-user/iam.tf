data "aws_iam_policy_document" "assume-role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals = {
      type        = "AWS"
      identifiers = ["${var.user_arn}"]
    }

    condition {
      test     = "Bool"
      variable = "aws:MultiFactorAuthPresent"
      values   = ["true"]
    }

    condition {
      test     = "IpAddress"
      variable = "aws:SourceIp"
      values   = ["${var.source_cidrs}"]
    }
  }
}

data "aws_iam_policy_document" "user-defaults" {
  statement {
    effect = "Allow"

    actions = [
      "autoscaling:Describe*",
      "cloudwatch:Describe*",
      "cloudwatch:Get*",
      "cloudwatch:List*",
      "logs:Get*",
      "logs:Describe*",
      "sns:Get*",
      "sns:List*",
      "eks:DescribeCluster*",
    ]

    resources = ["*"]
  }
}

resource "aws_iam_role" "user" {
  name               = "${var.cluster_name}-${var.role_prefix}-${var.user_name}"
  assume_role_policy = "${data.aws_iam_policy_document.assume-role.json}"
}

resource "aws_iam_policy" "user-defaults" {
  name   = "${var.cluster_name}-${var.user_name}-user-defaults"
  policy = "${data.aws_iam_policy_document.user-defaults.json}"
}

resource "aws_iam_policy_attachment" "user-defaults" {
  name       = "${var.cluster_name}-${var.user_name}-user-defaults-attachment"
  roles      = ["${aws_iam_role.user.name}"]
  policy_arn = "${aws_iam_policy.user-defaults.arn}"
}
