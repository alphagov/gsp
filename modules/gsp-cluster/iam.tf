resource "aws_iam_role" "aws-service-operator" {
  name               = "${var.cluster_name}-aws-service-operator"
  description        = "Role the AWS Service Operator assumes"
  assume_role_policy = "${data.aws_iam_policy_document.trust_kiam_server.json}"
}

data "aws_iam_policy_document" "aws-service-operator" {
  statement {
    actions = [
      "sqs:*",
      "sns:*",
      "cloudformation:*",
      "ecr:*",
      "dynamodb:*",
      "s3:*",
      "elasticache:*",
    ]

    resources = ["*"]
  }
}

resource "aws_iam_policy" "aws-service-operator" {
  name = "${var.cluster_name}-aws-service-operator"
  policy = "${data.aws_iam_policy_document.aws-service-operator.json}"
}

resource "aws_iam_role_policy_attachment" "aws-service-operator" {
  policy_arn = "${aws_iam_policy.aws-service-operator.arn}"
  role = "${aws_iam_role.aws-service-operator.name}"
}
