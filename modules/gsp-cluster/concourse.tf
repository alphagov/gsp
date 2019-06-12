resource "aws_iam_role" "concourse" {
  name        = "${var.cluster_name}-concourse"
  description = "Role the concourse process assumes"

  assume_role_policy = "${data.aws_iam_policy_document.assume-concourse-role.json}"
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
