data "aws_iam_policy_document" "user_data_policy_document" {
  statement {
    actions = [
      "s3:GetObject",
    ]

    resources = [
      "${data.aws_s3_bucket.user_data.arn}/user_data/${var.cluster_name}-*",
    ]
  }
}

data "aws_s3_bucket" "user_data" {
  bucket = "${var.user_data_bucket_name}"
}

resource "aws_iam_policy" "s3-user-data-policy" {
  name   = "${var.cluster_name}-s3-user-data-policy"
  policy = "${data.aws_iam_policy_document.user_data_policy_document.json}"
}
