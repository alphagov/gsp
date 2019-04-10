resource "aws_iam_role" "external-dns" {
  name        = "${var.cluster_name}-external-dns"
  description = "Role the external-dns process assumes"

  assume_role_policy = "${data.aws_iam_policy_document.assume-external-dns-role.json}"
}

resource "aws_iam_policy" "external-dns-route-53" {
  name        = "${var.cluster_name}-external-dns-route-53"
  description = "Policy for the external-dns Route 53 integration"

  policy = "${data.aws_iam_policy_document.external-dns-route-53.json}"
}

data "aws_iam_policy_document" "assume-external-dns-role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals = {
      type        = "AWS"
      identifiers = ["${aws_iam_role.kiam_server_role.arn}"]
    }
  }
}

data "aws_iam_policy_document" "external-dns-route-53" {
  statement {
    actions = [
      "route53:ChangeResourceRecordSets",
    ]

    resources = ["arn:aws:route53:::hostedzone/${data.aws_route53_zone.zone.zone_id}"]
  }

  statement {
    actions = [
      "route53:ListHostedZones",
      "route53:ListResourceRecordSets",
    ]

    resources = ["*"]
  }
}

resource "aws_iam_policy_attachment" "external-dns-route-53" {
  name       = "${var.cluster_name}-external-dns-route-53"
  roles      = ["${aws_iam_role.external-dns.name}"]
  policy_arn = "${aws_iam_policy.external-dns-route-53.arn}"
}

module "ingress-system" {
  enabled = 1
  source  = "../flux-release"

  namespace             = "ingress-system"
  chart_git             = "https://github.com/alphagov/gsp-ingress-system.git"
  chart_ref             = "master"
  cluster_name          = "${var.cluster_name}"
  cluster_domain        = "${var.cluster_name}.${var.dns_zone}"
  addons_dir            = "addons/${var.cluster_name}"
  permitted_roles_regex = "^${aws_iam_role.external-dns.name}$"

  extra_namespace_labels = <<EOF
    certmanager.k8s.io/disable-validation: "true"
EOF

  values = <<EOF
    webhook:
      enabled: false
    nginx-ingress:
      controller:
        service:
          annotations:
            external-dns.alpha.kubernetes.io/hostname: "${var.cluster_name}.${var.dns_zone}.,*.${var.cluster_name}.${var.dns_zone}."
            service.beta.kubernetes.io/aws-load-balancer-type: nlb
    external-dns:
      podAnnotations:
        iam.amazonaws.com/role: "${aws_iam_role.external-dns.name}"
      txtOwnerId: "${data.aws_route53_zone.zone.zone_id}"
EOF
}
