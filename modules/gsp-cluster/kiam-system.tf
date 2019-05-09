data "aws_iam_policy_document" "kiam_server_role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals = {
      type        = "AWS"
      identifiers = ["${module.k8s-cluster.kiam-server-node-instance-role-arn}"]
    }
  }
}

data "aws_iam_policy_document" "kiam_server_policy" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    resources = [
      "${aws_iam_role.cloudwatch_log_shipping_role.arn}",
      "${aws_iam_role.canary_role.arn}",
      "${aws_iam_role.flux-helm-operator.arn}",
    ]
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
