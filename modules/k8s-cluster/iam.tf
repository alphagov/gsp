data "aws_iam_policy_document" "eks-cluster-assume-role-policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["eks.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "ssm-minimal" {
  statement {
    actions = [
      "ssm:UpdateInstanceInformation",
      "ssmmessages:CreateControlChannel",
      "ssmmessages:CreateDataChannel",
      "ssmmessages:OpenControlChannel",
      "ssmmessages:OpenDataChannel",
    ]

    resources = ["*"]
  }

  statement {
    actions = [
      "s3:GetEncryptionConfiguration",
    ]

    resources = ["*"]
  }
}

data "aws_iam_policy_document" "ecr_access" {
  statement {
    actions = [
      "ecr:BatchCheckLayerAvailability",
      "ecr:BatchGetImage",
      "ecr:GetDownloadUrlForLayer",
      "ecr:GetAuthorizationToken",
    ]

    resources = ["*"]
  }
}

resource "aws_iam_role" "eks-cluster" {
  name               = "${var.cluster_name}-cluster"
  assume_role_policy = data.aws_iam_policy_document.eks-cluster-assume-role-policy.json
}

resource "aws_iam_role_policy_attachment" "eks-cluster-policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.eks-cluster.name
}

resource "aws_iam_role_policy_attachment" "eks-service-policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSServicePolicy"
  role       = aws_iam_role.eks-cluster.name
}

data "aws_arn" "worker-nodes-role" {
  arn = aws_cloudformation_stack.worker-nodes.outputs["NodeInstanceRole"]
}

data "aws_arn" "kiam-server-nodes-role" {
  arn = aws_cloudformation_stack.kiam-server-nodes.outputs["NodeInstanceRole"]
}

data "aws_arn" "ci-nodes-role" {
  arn = aws_cloudformation_stack.ci-nodes.outputs["NodeInstanceRole"]
}

resource "aws_iam_policy" "ssm-minimal" {
  name   = "${var.cluster_name}-ssm-minimal"
  policy = data.aws_iam_policy_document.ssm-minimal.json
}

resource "aws_iam_role_policy_attachment" "worker-nodes-ssm" {
  policy_arn = aws_iam_policy.ssm-minimal.arn
  role       = replace(data.aws_arn.worker-nodes-role.resource, "role/", "")
}

resource "aws_iam_role_policy_attachment" "kiam-nodes-ssm" {
  policy_arn = aws_iam_policy.ssm-minimal.arn
  role       = replace(data.aws_arn.kiam-server-nodes-role.resource, "role/", "")
}

resource "aws_iam_role_policy_attachment" "ci-nodes-ssm" {
  policy_arn = aws_iam_policy.ssm-minimal.arn
  role       = replace(data.aws_arn.ci-nodes-role.resource, "role/", "")
}

resource "aws_iam_policy" "ecr_access" {
  name   = "${var.cluster_name}-ecr-access"
  policy = data.aws_iam_policy_document.ecr_access.json
}

resource "aws_iam_role_policy_attachment" "worker_nodes_ecr_access" {
  policy_arn = aws_iam_policy.ecr_access.arn
  role       = replace(data.aws_arn.worker-nodes-role.resource, "role/", "")
}

resource "aws_iam_role_policy_attachment" "ci_nodes_ecr_access" {
  policy_arn = aws_iam_policy.ecr_access.arn
  role       = replace(data.aws_arn.ci-nodes-role.resource, "role/", "")
}
