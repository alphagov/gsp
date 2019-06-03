resource "aws_iam_role" "concourse" {
  name        = "${var.cluster_name}-concourse"
  description = "Role the concourse process assumes"

  assume_role_policy = "${data.aws_iam_policy_document.assume-concourse-role.json}"
}

resource "aws_iam_policy" "concourse-code-commit" {
  name        = "${var.cluster_name}-concourse-code-commit"
  description = "Policy for the concourse to access code commit"

  policy = "${data.aws_iam_policy_document.concourse-code-commit.json}"
}

data "aws_iam_policy_document" "assume-concourse-role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals = {
      type        = "AWS"
      identifiers = ["${aws_iam_role.kiam_server_role.arn}"]
    }
  }
}

data "aws_iam_policy_document" "concourse-code-commit" {
  statement {
    actions = [
      "codecommit:GitPull",
    ]

    resources = [
      "${aws_codecommit_repository.canary.arn}",
    ]
  }
}

resource "aws_iam_policy_attachment" "concourse-code-commit" {
  name       = "${var.cluster_name}-concourse-code-commit"
  roles      = ["${aws_iam_role.concourse.name}"]
  policy_arn = "${aws_iam_policy.concourse-code-commit.arn}"
}
