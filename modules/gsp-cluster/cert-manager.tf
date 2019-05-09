resource "aws_iam_role" "cert-manager" {
  name        = "${var.cluster_name}-cert-manager"
  description = "Role the cert-manager assumes"

  assume_role_policy = "${data.aws_iam_policy_document.assume-cert-manager-role.json}"
}

resource "aws_iam_policy" "cert-manager-route53" {
  name        = "${var.cluster_name}-cert-manager-route53"
  description = "Policy for the cert-manager to access route53"

  policy = "${data.aws_iam_policy_document.cert-manager-route53.json}"
}

data "aws_iam_policy_document" "assume-cert-manager-role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals = {
      type        = "AWS"
      identifiers = ["${aws_iam_role.kiam_server_role.arn}"]
    }
  }
}

data "aws_iam_policy_document" "cert-manager-route53" {
  statement {
    actions = [
      "route53:GetChange",
    ]

    resources = [
      "arn:aws:route53:::change/*",
    ]
  }
  
  statement {
    actions = [
      "route53:ChangeResourceRecordSets",
    ]

    resources = [
      "arn:aws:route53:::hostedzone/${var.cluster_domain_id}",
    ]
  }
  
  statement {
    actions = [
      "route53:ListHostedZonesByName",
    ]

    resources = [
      "*",
    ]
  }
}

resource "aws_iam_policy_attachment" "cert-manager-route53" {
  name       = "${var.cluster_name}-cert-manager-route53"
  roles      = ["${aws_iam_role.cert-manager.name}"]
  policy_arn = "${aws_iam_policy.cert-manager-route53.arn}"
}
