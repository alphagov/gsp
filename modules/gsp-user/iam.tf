data "aws_iam_policy_document" "assume-role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "AWS"
      identifiers = [var.user_arn]
    }

    condition {
      test     = "Bool"
      variable = "aws:MultiFactorAuthPresent"
      values   = ["true"]
    }

    condition {
      test     = "IpAddress"
      variable = "aws:SourceIp"
      values   = var.source_cidrs
    }
  }
}

data "aws_iam_policy_document" "user-defaults" {
  statement {
    effect = "Allow"

    actions = [
      "eks:DescribeCluster*",
    ]

    resources = ["*"]
  }
}

resource "aws_iam_role" "user" {
  name               = "${var.cluster_name}-${var.role_prefix}-${var.user_name}"
  assume_role_policy = data.aws_iam_policy_document.assume-role.json
}

resource "aws_iam_policy" "user-defaults" {
  name   = "${var.cluster_name}-${var.user_name}-user-defaults"
  policy = data.aws_iam_policy_document.user-defaults.json
}

resource "aws_iam_role_policy_attachment" "user-defaults" {
  role       = aws_iam_role.user.name
  policy_arn = aws_iam_policy.user-defaults.arn
}

resource "aws_iam_role_policy_attachment" "user-defaults-cloudwatch" {
  role       = aws_iam_role.user.name
  policy_arn = "arn:aws:iam::aws:policy/CloudWatchReadOnlyAccess"
}

resource "aws_iam_role_policy_attachment" "user-defaults-cloudformation" {
  role       = aws_iam_role.user.name
  policy_arn = "arn:aws:iam::aws:policy/AWSCloudFormationReadOnlyAccess"
}

resource "aws_iam_role_policy_attachment" "user-defaults-aws-support" {
  role       = aws_iam_role.user.name
  policy_arn = "arn:aws:iam::aws:policy/AWSSupportAccess"
}

resource "aws_iam_role_policy_attachment" "user-defaults-iam" {
  role       = aws_iam_role.user.name
  policy_arn = "arn:aws:iam::aws:policy/IAMReadOnlyAccess"
}

resource "aws_iam_role_policy_attachment" "user-defaults-rds" {
  role       = aws_iam_role.user.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonRDSReadOnlyAccess"
}

resource "aws_iam_role_policy_attachment" "user-defaults-load-balancing-read-only" {
  role       = aws_iam_role.user.name
  policy_arn = "arn:aws:iam::aws:policy/ElasticLoadBalancingReadOnly"
}

resource "aws_iam_role_policy_attachment" "user-defaults-view-only" {
  role       = aws_iam_role.user.name
  policy_arn = "arn:aws:iam::aws:policy/job-function/ViewOnlyAccess"
}
