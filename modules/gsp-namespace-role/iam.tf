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

data "aws_iam_policy_document" "namespace-sqs" {
  statement {
    effect = "Allow"

    actions = [
      "sqs:SendMessage",
      "sqs:ReceiveMessage",
      "sqs:DeleteMessage"
    ]

    resources = ["arn:aws:sqs:*:${var.account_id}:${var.cluster_name}-sqsqueue-${var.namespace_name}-*"]
  }
}

resource "aws_iam_role" "namespace-sqs" {
  name               = "${var.cluster_name}-namespace-${var.namespace_name}-sqs"
  assume_role_policy = "${data.aws_iam_policy_document.assume-role.json}"
  path               = "/gsp/${var.cluster_name}/namespaceroles/sqs/"
}

resource "aws_iam_policy" "namespace-sqs" {
  name   = "${var.cluster_name}-namespace-${var.namespace_name}-sqs"
  policy = "${data.aws_iam_policy_document.namespace-sqs.json}"
}

resource "aws_iam_role_policy_attachment" "namespace-sqs" {
  role       = "${aws_iam_role.namespace.name}"
  policy_arn = "${aws_iam_policy.namespace-sqs.arn}"
}
