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

resource "aws_iam_role" "namespace" {
  name               = "${var.cluster_name}-namespace-${var.namespace_name}-${var.role_name}"
  assume_role_policy = "${data.aws_iam_policy_document.assume-role.json}"
  path               = "/gsp/${var.cluster_name}/namespaceroles/"
}
