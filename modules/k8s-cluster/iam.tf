data "aws_iam_policy_document" "controller_role_doc" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "controller_role" {
  name               = "${var.cluster_name}-controller-instance-role"
  assume_role_policy = "${data.aws_iam_policy_document.controller_role_doc.json}"
}

data "aws_iam_policy_document" "controller_policy_doc" {
  statement {
    actions = [
      "elasticloadbalancing:DescribeLoadBalancers" # https://github.com/kubernetes/kubernetes/issues/47733
    ]
    resources = ["*"]
  }

  statement {
    actions = [
      "sts:AssumeRole",
    ]

    resources = ["*"] # This allows kiam to assume any role. We rely on trust relationships from the other side to ensure that it can't assume everything.
  }
}

resource "aws_iam_role_policy_attachment" "controller-policy-attachment" {
  role       = "${aws_iam_role.controller_role.id}"
  policy_arn = "${aws_iam_policy.controller-policy.arn}"
}

resource "aws_iam_policy" "controller-policy" {
  name   = "${var.cluster_name}-controller-instance-policy"
  policy = "${data.aws_iam_policy_document.controller_policy_doc.json}"
}

data "aws_iam_policy_document" "controller_persistent_volume_claim" {
  # https://docs.docker.com/ee/ucp/kubernetes/storage/configure-aws-storage/
  statement {
    actions = [
      "ec2:DescribeInstances",
      "ec2:DescribeVolumes",
      "ec2:CreateVolume",
    ]

    resources = ["*"]
  }

  statement {
    actions = [
      "ec2:CreateTags",
    ]

    resources = ["arn:aws:ec2:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:volume/*"]

    condition {
      "test" = "StringEquals"
      "variable" = "aws:RequestTag/KubernetesCluster"
      "values" = [
        "${var.cluster_name}",
      ]
    }

    condition {
      test = "ForAnyValue:StringEquals"
      variable = "aws:TagKeys"
      values = [
        "KubernetesCluster"
      ]
    }
  }

  statement {
    actions = [
      "ec2:DetachVolume",
      "ec2:AttachVolume",
      "ec2:DeleteVolume",
    ]

    resources = ["*"]

    condition {
      "test" = "StringEquals"
      "variable" = "ec2:ResourceTag/KubernetesCluster"
      "values" = [
        "${var.cluster_name}",
      ]
    }
  }
}

resource "aws_iam_policy" "controller-persistent-volume-claim-policy" {
  name   = "${var.cluster_name}-controller-persistent-volume-claim-policy"
  policy = "${data.aws_iam_policy_document.controller_persistent_volume_claim.json}"
}

resource "aws_iam_role_policy_attachment" "controller_persistent_volume_claim_attachement" {
  role       = "${aws_iam_role.controller_role.id}"
  policy_arn = "${aws_iam_policy.controller-persistent-volume-claim-policy.arn}"
}

resource "aws_iam_role_policy_attachment" "controller-ssm-policy-attachment" {
  role       = "${aws_iam_role.controller_role.id}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforSSM"
}

resource "aws_iam_role_policy_attachment" "controller-s3-user-data-policy-attachment" {
  role       = "${aws_iam_role.controller_role.id}"
  policy_arn = "${var.s3_user_data_policy_arn}"
}

resource "aws_iam_instance_profile" "controller_profile" {
  name = "${var.cluster_name}-controller-instance-role"
  role = "${aws_iam_role.controller_role.name}"
}

data "aws_iam_policy_document" "worker_role_doc" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "worker_role" {
  name               = "${var.cluster_name}-worker-instance-role"
  assume_role_policy = "${data.aws_iam_policy_document.worker_role_doc.json}"
}

data "aws_iam_policy_document" "worker_policy_doc" {
  statement {
    actions = [
      "ec2:DescribeInstances",
      "ec2:DescribeRegions",
    ]

    resources = ["*"]
  }

  statement {
    actions   = ["sts:AssumeRole"]
    resources = ["*"]
  }
}

resource "aws_iam_policy" "worker-policy" {
  name   = "${var.cluster_name}-worker-instance-policy"
  policy = "${data.aws_iam_policy_document.worker_policy_doc.json}"
}

resource "aws_iam_role_policy_attachment" "worker-policy-attachment" {
  role       = "${aws_iam_role.worker_role.id}"
  policy_arn = "${aws_iam_policy.worker-policy.arn}"
}

resource "aws_iam_role_policy_attachment" "worker-ssm-policy-attachment" {
  role       = "${aws_iam_role.worker_role.id}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforSSM"
}

resource "aws_iam_role_policy_attachment" "worker-s3-user-data-policy-attachment" {
  role       = "${aws_iam_role.worker_role.id}"
  policy_arn = "${var.s3_user_data_policy_arn}"
}

resource "aws_iam_instance_profile" "worker_profile" {
  name = "${var.cluster_name}-worker-instance-role"
  role = "${aws_iam_role.worker_role.name}"
}
