data "aws_iam_policy_document" "assume-harbor" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals = {
      type        = "AWS"
      identifiers = ["${var.harbor_role_asumer_arn}"]
    }
  }
}

resource "aws_iam_role" "harbor" {
  count = "${var.enabled == 0 ? 0 : 1}"

  name        = "${var.cluster_name}-harbor"
  description = "Role the harbor process assumes"

  assume_role_policy = "${data.aws_iam_policy_document.assume-harbor.json}"
}

data "aws_iam_policy_document" "harbor-s3" {
  statement {
    actions = [
      "s3:*",
    ]

    resources = [
      "${element(concat(aws_s3_bucket.ci-system-harbor-registry-storage.*.arn, list("")), 0)}",
      "${element(concat(aws_s3_bucket.ci-system-harbor-registry-storage.*.arn, list("")), 0)}/*",
    ]
  }
}

resource "aws_iam_policy" "harbor-s3" {
  count = "${var.enabled == 0 ? 0 : 1}"

  name        = "${var.cluster_name}-harbor-s3"
  description = "Policy for the harbor s3 access"

  policy = "${data.aws_iam_policy_document.harbor-s3.json}"
}

resource "aws_iam_policy_attachment" "harbor-s3" {
  count = "${var.enabled == 0 ? 0 : 1}"

  name       = "${var.cluster_name}-harbor-s3"
  roles      = ["${element(aws_iam_role.harbor.*.name, count.index)}"]
  policy_arn = "${element(aws_iam_policy.harbor-s3.*.arn, count.index)}"
}
