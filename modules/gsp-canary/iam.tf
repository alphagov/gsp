resource "aws_iam_role" "canary_role" {
  name        = "${var.cluster_name}-canary"
  description = "Role the gsp-canary process assumes"

  assume_role_policy = "${data.aws_iam_policy_document.assume_canary_role.json}"
}

resource "aws_iam_policy" "canary_code_commit" {
  name        = "${var.cluster_name}-canary-code-commit"
  description = "Policy for the gsp-canary code commit"

  policy = "${data.aws_iam_policy_document.canary_code_commit.json}"
}

data "aws_iam_policy_document" "assume_canary_role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals = {
      type        = "AWS"
      identifiers = ["${var.canary_role_assumer_arn}"]
    }
  }
}

data "aws_iam_policy_document" "canary_code_commit" {
  statement {
    actions = [
      "codecommit:GitPull",
      "codecommit:GitPush",
    ]

    resources = ["${aws_codecommit_repository.canary.arn}"]
  }
}

resource "aws_iam_policy_attachment" "canary_code_commit" {
  name       = "${var.cluster_name}-canary-code-commit"
  roles      = ["${aws_iam_role.canary_role.name}"]
  policy_arn = "${aws_iam_policy.canary_code_commit.arn}"
}
