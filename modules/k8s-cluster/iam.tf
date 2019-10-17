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

resource "aws_iam_role" "eks-cluster" {
  name               = "${var.cluster_name}-cluster"
  assume_role_policy = "${data.aws_iam_policy_document.eks-cluster-assume-role-policy.json}"
}

resource "aws_iam_role_policy_attachment" "eks-cluster-policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = "${aws_iam_role.eks-cluster.name}"
}

resource "aws_iam_role_policy_attachment" "eks-service-policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSServicePolicy"
  role       = "${aws_iam_role.eks-cluster.name}"
}

data "aws_iam_policy_document" "worker-nodes-assume-role-policy" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals = {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "worker-nodes-role" {
  name               = "${var.cluster_name}-worker-nodes-role"
  assume_role_policy = "${data.aws_iam_policy_document.worker-nodes-assume-role-policy.json}"
}

resource "aws_iam_instance_profile" "worker-nodes-profile" {
  name = "${var.cluster_name}-worker-nodes-profile"
  role = "${aws_iam_role.worker-nodes-role.name}"
}

resource "aws_iam_role_policy_attachment" "worker-nodes-eks-worker-policy-attachment" {
  role       = "${aws_iam_role.worker-nodes-role.name}"
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
}

resource "aws_iam_role_policy_attachment" "worker-nodes-eks-cni-policy-attachment" {
  role       = "${aws_iam_role.worker-nodes-role.name}"
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
}

resource "aws_iam_role_policy_attachment" "worker-nodes-ecr-ro-policy-attachment" {
  role       = "${aws_iam_role.worker-nodes-role.name}"
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}

data "aws_arn" "kiam-server-nodes-role" {
  arn = "${aws_cloudformation_stack.kiam-server-nodes.outputs["NodeInstanceRole"]}"
}

data "aws_arn" "ci-nodes-role" {
  arn = "${aws_cloudformation_stack.ci-nodes.outputs["NodeInstanceRole"]}"
}

resource "aws_iam_policy" "ssm-minimal" {
  name = "${var.cluster_name}-ssm-minimal"
  policy = "${data.aws_iam_policy_document.ssm-minimal.json}"
}

resource "aws_iam_role_policy_attachment" "worker-nodes-ssm" {
  policy_arn = "${aws_iam_policy.ssm-minimal.arn}"
  role = "${aws_iam_role.worker-nodes-role.name}"
}

resource "aws_iam_role_policy_attachment" "kiam-nodes-ssm" {
  policy_arn = "${aws_iam_policy.ssm-minimal.arn}"
  role = "${replace(data.aws_arn.kiam-server-nodes-role.resource, "role/", "")}"
}

resource "aws_iam_role_policy_attachment" "ci-nodes-ssm" {
  policy_arn = "${aws_iam_policy.ssm-minimal.arn}"
  role = "${replace(data.aws_arn.ci-nodes-role.resource, "role/", "")}"
}
