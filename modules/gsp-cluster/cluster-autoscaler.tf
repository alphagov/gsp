resource "aws_iam_role" "cluster_autoscaler" {
  name = "cluster-autoscaler"

  assume_role_policy = "${data.aws_iam_policy_document.trust_kiam_server.json}"
}

data "aws_iam_policy_document" "cluster_autoscaler_policy" {
  statement {
    effect = "Allow"

    actions = [
      "autoscaling:DescribeAutoScalingGroups",
      "autoscaling:DescribeAutoScalingInstances",
      "autoscaling:DescribeLaunchConfigurations",
      "autoscaling:DescribeTags",
      "ec2:DescribeLaunchTemplateVersions",
    ]

    resources = ["*"]
  }

  statement {
    effect = "Allow"

    actions = [
      "autoscaling:SetDesiredCapacity",
      "autoscaling:TerminateInstanceInAutoScalingGroup",
    ]

    condition = {
      test     = "Null"
      variable = "autoscaling:ResourceTag/k8s.io/cluster-autoscaler/${var.cluster_name}"
      values   = ["false"]
    }

    resources = ["*"]
  }
}

resource "aws_iam_policy" "cluster-autoscaler" {
  name        = "${var.cluster_name}-cluster-autoscaler"
  description = "Policy for the cluster autoscaler"
  policy      = "${data.aws_iam_policy_document.cluster_autoscaler_policy.json}"
}

resource "aws_iam_policy_attachment" "cluster-autoscaler" {
  name       = "${var.cluster_name}-cluster-autoscaler"
  roles      = ["${aws_iam_role.cluster_autoscaler.name}"]
  policy_arn = "${aws_iam_policy.cluster-autoscaler.arn}"
}

resource "aws_iam_policy_attachment" "cluster-autoscaler-mgmt" {
  name       = "${var.cluster_name}-cluster-autoscaler-mgmt"
  roles      = ["${module.k8s-cluster.kiam-server-node-instance-role-name}"]
  policy_arn = "${aws_iam_policy.cluster-autoscaler.arn}"
}
