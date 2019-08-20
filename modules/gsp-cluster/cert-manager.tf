data "aws_iam_policy_document" "cert_manager" {
  statement {
    effect = "Allow"

    actions = [
      "route53:GetHostedZone",
      "route53:ChangeResourceRecordSets",
      "route53:ListResourceRecordSets",
    ]

    resources = [
      "arn:aws:route53:::hostedzone/${var.cluster_domain_id}"
    ]
  }

  statement {
    effect = "Allow"

    actions = [
      "route53:ListHostedZones",
    ]

    resources = [
      "*"
    ]
  }
}

resource "aws_iam_policy" "cert_manager" {
  name         = "${var.cluster_name}_cert_manager"
  description = "Allow cert-manager to use the DNS01 challenge"

  policy = "${data.aws_iam_policy_document.cert_manager.json}"
}

resource "aws_iam_role" "cert_manager" {
  name = "${var.cluster_name}_cert_manager"

  assume_role_policy = "${data.aws_iam_policy_document.trust_kiam_server.json}"
}

resource "aws_iam_policy_attachment" "cert_manager" {
  name       = "${var.cluster_name}_cert_manager"
  roles      = [
    "${aws_iam_role.cert_manager.name}",
  ]
  policy_arn = "${aws_iam_policy.cert_manager.arn}"
}
