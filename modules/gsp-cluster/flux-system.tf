resource "aws_iam_role" "flux-helm-operator" {
  name        = "${var.cluster_name}-flux-helm-operator"
  description = "Role the flux helm operator process assumes"

  assume_role_policy = "${data.aws_iam_policy_document.assume-flux-helm-operator-role.json}"
}

resource "aws_iam_policy" "flux-helm-operator-code-commit" {
  name        = "${var.cluster_name}-flux-helm-operator-code-commit"
  description = "Policy for the flux helm operator to access code commit"

  policy = "${data.aws_iam_policy_document.flux-helm-operator-code-commit.json}"
}

data "aws_iam_policy_document" "assume-flux-helm-operator-role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals = {
      type        = "AWS"
      identifiers = ["${aws_iam_role.kiam_server_role.arn}"]
    }
  }
}

data "aws_iam_policy_document" "flux-helm-operator-code-commit" {
  statement {
    actions = [
      "codecommit:GitPull",
    ]

    resources = [
      "${aws_codecommit_repository.canary.arn}",
    ]
  }
}

resource "aws_iam_policy_attachment" "flux-helm-operator-code-commit" {
  name       = "${var.cluster_name}-flux-helm-operator-code-commit"
  roles      = ["${aws_iam_role.flux-helm-operator.name}"]
  policy_arn = "${aws_iam_policy.flux-helm-operator-code-commit.arn}"
}
