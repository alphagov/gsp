resource "aws_iam_role_policy_attachment" "namespace-sqs" {
  role       = "${var.cluster_name}-namespace-${var.namespace_name}"
  policy_arn = "arn:aws:iam::${var.account_id}:role/${var.cluster_name}-namespace-${var.namespace_name}-${var.role_name}"
}
