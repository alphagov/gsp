data "aws_iam_policy_document" "user_data_policy_document" {
  statement {
    actions = [
      "s3:GetObject",
    ]

    resources = [
      "${data.aws_s3_bucket.user_data.arn}/*",
    ]

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }

    condition {
      test     = "StringEquals"
      variable = "aws:sourceVpc"

      values = ["${aws_vpc.network.id}"]
    }
  }
}

data "aws_s3_bucket" "user_data" {
  bucket = "${var.user_data_bucket_name}"
}

resource "aws_s3_bucket_policy" "user_data_bucket_policy" {
  bucket = "${data.aws_s3_bucket.user_data.id}"
  policy = "${data.aws_iam_policy_document.user_data_policy_document.json}"
}
