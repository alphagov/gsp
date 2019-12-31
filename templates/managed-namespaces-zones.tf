{{- $allNamespaces := prepend (dict "name" "gsp-system" "ingress" (dict "enabled" true)) (datasource "config").namespaces }}
{{- range $namespace := $allNamespaces }}
{{- if (has $namespace "ingress") }}
{{- if (has $namespace.ingress "enabled") }}
{{- if $namespace.ingress.enabled }}

resource "aws_route53_zone" "{{ $namespace.name }}" {
  name          = "{{ $namespace.name }}.${var.cluster_domain}"
  force_destroy = true
}

resource "aws_route53_record" "{{ $namespace.name }}-ns" {
  zone_id  = module.gsp-domain.zone_id
  name     = "{{ $namespace.name }}.${var.cluster_domain}"
  type     = "NS"
  ttl      = "30"

  records = [
    aws_route53_zone.{{ $namespace.name }}.name_servers[0],
    aws_route53_zone.{{ $namespace.name }}.name_servers[1],
    aws_route53_zone.{{ $namespace.name }}.name_servers[2],
    aws_route53_zone.{{ $namespace.name }}.name_servers[3],
  ]
}

data "aws_iam_policy_document" "trust_{{ $namespace.name }}_external_dns" {
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
      values = ["system:serviceaccount:{{ $namespace.name }}:gsp-external-dns"]
    }
  }
}

data "aws_iam_policy_document" "{{ $namespace.name }}-external-dns" {
  statement {
    effect = "Allow"

    actions = [
      "route53:ChangeResourceRecordSets",
    ]

    resources = [
      "arn:aws:route53:::hostedzone/${aws_route53_zone.{{ $namespace.name }}.zone_id}",
{{- if eq "gsp-system" $namespace.name }}
      "arn:aws:route53:::hostedzone/${module.gsp-domain.zone_id}",
{{- end }}
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

resource "aws_iam_policy" "{{ $namespace.name }}-external-dns" {
  name        = "${var.cluster_name}-{{ $namespace.name }}-external-dns"
  description = "Allow external-dns to do its job in namespace {{ $namespace.name }}"

  policy = data.aws_iam_policy_document.{{ $namespace.name }}-external-dns.json
}

resource "aws_iam_role" "{{ $namespace.name }}-external-dns" {
  name = "${var.cluster_name}-{{ $namespace.name }}-external-dns"

  assume_role_policy = data.aws_iam_policy_document.trust_{{ $namespace.name }}_external_dns.json
}

resource "aws_iam_policy_attachment" "{{ $namespace.name }}-external-dns" {
  name = "${var.cluster_name}-{{ $namespace.name }}-external-dns"
  roles = [
    aws_iam_role.{{ $namespace.name }}-external-dns.name,
  ]
  policy_arn = aws_iam_policy.{{ $namespace.name }}-external-dns.arn
}

{{- end }}
{{- end }}
{{- end }}
{{- end }}

locals {
  external-dns-namespace-zones = [
{{- range $namespace := $allNamespaces }}
{{- if (has $namespace "ingress") }}
{{- if (has $namespace.ingress "enabled") }}
{{- if $namespace.ingress.enabled }}
    {
      namespace = "{{ $namespace.name }}",
      zoneId = aws_route53_zone.{{ $namespace.name }}.zone_id,
      roleArn = aws_iam_role.{{ $namespace.name }}-external-dns.arn
    },
{{- end }}
{{- end }}
{{- end }}
{{- end }}
  ]
  cluster_zone_ids = [
    module.gsp-domain.zone_id,
{{- range $namespace := $allNamespaces }}
{{- if (has $namespace "ingress") }}
{{- if (has $namespace.ingress "enabled") }}
{{- if $namespace.ingress.enabled }}
    aws_route53_zone.{{ $namespace.name }}.zone_id,
{{- end }}
{{- end }}
{{- end }}
{{- end }}
  ]
}
