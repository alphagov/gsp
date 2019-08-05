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

resource "aws_iam_policy" "namespace-sqs" {
  name   = "${var.cluster_name}-namespace-${var.namespace_name}-sqs"
  policy = "${data.aws_iam_policy_document.namespace-sqs.json}"
}
