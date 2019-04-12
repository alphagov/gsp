data "aws_iam_policy_document" "eks-cluster-assume-role-policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["eks.amazonaws.com"]
    }
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

data "aws_arn" "worker-nodes-role" {
  arn = "${aws_cloudformation_stack.worker-nodes.outputs["NodeInstanceRole"]}"
}

resource "aws_iam_role_policy_attachment" "worker-nodes-ssm" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforSSM"
  role = "${replace(data.aws_arn.worker-nodes-role.resource, "role/", "")}"
}

data "aws_arn" "kiam-server-nodes-role" {
  arn = "${aws_cloudformation_stack.kiam-server-nodes.outputs["NodeInstanceRole"]}"
}

resource "aws_iam_role_policy_attachment" "kiam-server-nodes-ssm" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforSSM"
  role = "${replace(data.aws_arn.kiam-server-nodes-role.resource, "role/", "")}"
}

data "aws_arn" "ci-nodes-role" {
  arn = "${aws_cloudformation_stack.ci-nodes.outputs["NodeInstanceRole"]}"
}

resource "aws_iam_role_policy_attachment" "ci-nodes-ssm" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforSSM"
  role = "${replace(data.aws_arn.ci-nodes-role.resource, "role/", "")}"
}
