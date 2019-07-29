data "aws_iam_policy_document" "assume-role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals = {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${var.account_id}:role/${var.cluster_name}_kiam_server"]
    }
  }
}

data "aws_iam_policy_document" "namespace-defaults" {
  statement {
    effect = "Allow"

    actions = [
      "sqs:SendMessage",
      "sqs:ReceiveMessage"
    ]

    resources = ["arn:aws:sqs:*:${var.account_id}:${var.cluster_name}-sqsqueue-${var.namespace_name}-*"]
  }
}

resource "aws_iam_role" "namespace" {
  name               = "${var.cluster_name}-namespace-${var.namespace_name}"
  assume_role_policy = "${data.aws_iam_policy_document.assume-role.json}"
  path               = "/gsp/${var.cluster_name}/namespaceroles/"
}

resource "aws_iam_policy" "namespace-defaults" {
  name   = "${var.cluster_name}-namespace-${var.namespace_name}-defaults"
  policy = "${data.aws_iam_policy_document.namespace-defaults.json}"
}

resource "aws_iam_role_policy_attachment" "namespace-defaults" {
  role       = "${aws_iam_role.namespace.name}"
  policy_arn = "${aws_iam_policy.namespace-defaults.arn}"
}