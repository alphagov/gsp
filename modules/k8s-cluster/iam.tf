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
      "ec2:*",
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

resource "aws_iam_policy" "controller-policy" {
  name   = "${var.cluster_name}-controller-instance-policy"
  policy = "${data.aws_iam_policy_document.controller_policy_doc.json}"
}

resource "aws_iam_role_policy_attachment" "controller-policy-attachment" {
  role       = "${aws_iam_role.controller_role.id}"
  policy_arn = "${aws_iam_policy.controller-policy.arn}"
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
