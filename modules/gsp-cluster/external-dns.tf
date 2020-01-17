data "aws_iam_policy_document" "trust_external_dns" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRoleWithWebIdentity"]

    principals {
      type        = "Federated"
      identifiers = [module.k8s-cluster.oidc_provider_arn]
    }

    condition {
      test = "StringEquals"
      variable = "${replace(module.k8s-cluster.oidc_provider_url, "https://", "")}:sub"
      values = ["system:serviceaccount:gsp-system:gsp-external-dns"]
    }
  }
}

data "aws_iam_policy_document" "external_dns" {
  statement {
    effect = "Allow"

    actions = [
      "route53:ChangeResourceRecordSets",
    ]

    resources = [
      "arn:aws:route53:::hostedzone/${var.cluster_domain_id}",
    ]
  }

  statement {
    effect = "Allow"

    actions = [
      "route53:ListHostedZones",
      "route53:ListResourceRecordSets",
    ]

    resources = [
      "*",
    ]
  }
}

resource "aws_iam_policy" "external_dns" {
  name        = "${var.cluster_name}_external_dns"
  description = "Allow external-dns to do its job"

  policy = data.aws_iam_policy_document.external_dns.json
}

resource "aws_iam_role" "external_dns" {
  name = "${var.cluster_name}_external_dns"

  assume_role_policy = data.aws_iam_policy_document.trust_external_dns.json
}

resource "aws_iam_policy_attachment" "external_dns" {
  name = "${var.cluster_name}_external_dns"
  roles = [
    aws_iam_role.external_dns.name,
  ]
  policy_arn = aws_iam_policy.external_dns.arn
}

