resource "aws_iam_role" "lambda_log_forwarder" {
  count              = var.enabled == "0" ? 0 : 1
  name               = "${var.cluster_name}_${var.name}_lambda_log_forwarder"
  assume_role_policy = data.aws_iam_policy_document.lambda_log_forwarder_assume_role_policy.json
}

resource "aws_iam_policy_attachment" "lambda_log_forwarder" {
  count      = var.enabled == "0" ? 0 : 1
  name       = "${var.cluster_name}_${var.name}_lambda_log_forwarder_attachment"
  roles      = ["${aws_iam_role.lambda_log_forwarder[0].name}"]
  policy_arn = "${aws_iam_policy.lambda_log_forwarder[0].arn}"
}

resource "aws_iam_policy" "lambda_log_forwarder" {
  count       = var.enabled == "0" ? 0 : 1
  name        = "${var.cluster_name}_${var.name}_lambda_log_forwarder"
  description = "Policy for Lambda log forwarding function"
  policy      = data.aws_iam_policy_document.lambda_log_forwarder.json
}

data "aws_iam_policy_document" "lambda_log_forwarder" {
  statement {
    effect = "Allow"

    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = ["*"]
  }
}

data "aws_iam_policy_document" "lambda_log_forwarder_assume_role_policy" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

