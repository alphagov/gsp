data "aws_iam_policy_document" "trust_grafana" {
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
      values = ["system:serviceaccount:gsp-system:gsp-grafana"]
    }
  }
}

resource "aws_iam_role" "grafana" {
  name               = "${var.cluster_name}-grafana"
  description        = "Role the Grafana process assumes"
  assume_role_policy = data.aws_iam_policy_document.trust_grafana.json
}

data "aws_iam_policy_document" "grafana_cloudwatch" {
  statement {
    effect = "Allow"

    actions = [
      "cloudwatch:DescribeAlarmsForMetric",
      "cloudwatch:ListMetrics",
      "cloudwatch:GetMetricStatistics",
      "cloudwatch:GetMetricData",
    ]

    resources = ["*"]
  }

  statement {
    effect = "Allow"

    actions = [
      "ec2:DescribeTags",
      "ec2:DescribeInstances",
      "ec2:DescribeRegions",
    ]

    resources = ["*"]
  }

  statement {
    effect = "Allow"

    actions = [
      "tag:GetResources",
    ]

    resources = ["*"]
  }
}

resource "aws_iam_role_policy" "grafana" {
  name   = "${var.cluster_name}-grafana"
  role   = aws_iam_role.grafana.id
  policy = data.aws_iam_policy_document.grafana_cloudwatch.json
}

resource "random_password" "grafana_default_admin_password" {
  length  = 40
  special = false
}

