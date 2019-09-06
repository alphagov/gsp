resource "aws_iam_role" "sentry" {
  name               = "${var.cluster_name}-sentry"
  description        = "Role the sentry process assumes"
  assume_role_policy = "${data.aws_iam_policy_document.trust_kiam_server.json}"
}

data "aws_iam_policy_document" "sentry-s3" {
  statement {
    actions = [
      "s3:*",
    ]

    resources = [
      "${aws_s3_bucket.gsp-system-sentry-storage.arn}",
      "${aws_s3_bucket.gsp-system-sentry-storage.arn}/*",
    ]
  }
}

resource "aws_iam_policy" "sentry-s3" {
  name        = "${var.cluster_name}-sentry-s3"
  description = "Policy for the sentry s3 access"
  policy      = "${data.aws_iam_policy_document.sentry-s3.json}"
}

resource "aws_iam_policy_attachment" "sentry-s3" {
  name       = "${var.cluster_name}-sentry-s3"
  roles      = ["${aws_iam_role.sentry.name}"]
  policy_arn = "${aws_iam_policy.sentry-s3.arn}"
}

resource "aws_s3_bucket" "gsp-system-sentry-storage" {
  bucket = "sentry-${var.cluster_name}-${var.account_name}"
  acl    = "private"

  force_destroy = true # NEED TO VALIDATE!!!

  tags = {
    Name = "Sentry's persistant storage"
  }
}
