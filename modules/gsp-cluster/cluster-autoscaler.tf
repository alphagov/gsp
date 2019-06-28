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
      "autoscaling:SetDesiredCapacity",
      "autoscaling:TerminateInstanceInAutoScalingGroup",
    ]

    # TODO: can we restrict this to only the ASGs we care about?  we
    # can't construct the ARN because it has a UUID in it so we'd have
    # to fish it out of the cloudformation results somehow, or use a
    # `data.aws_autoscaling_group[s]` resource to identify it
    resources = ["*"]
  }
}
