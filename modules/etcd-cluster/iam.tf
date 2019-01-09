data "aws_iam_policy_document" "etcd_role_doc" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "etcd_role" {
  name               = "${var.cluster_name}-etcd-instance-role"
  assume_role_policy = "${data.aws_iam_policy_document.etcd_role_doc.json}"
}

resource "aws_iam_role_policy_attachment" "etcd-ssm-policy-attachment" {
  role       = "${aws_iam_role.etcd_role.id}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforSSM"
}

resource "aws_iam_role_policy_attachment" "etcd-s3-user-data-policy-attachment" {
  role       = "${aws_iam_role.etcd_role.id}"
  policy_arn = "${var.s3_user_data_policy_arn}"
}

resource "aws_iam_instance_profile" "etcd_profile" {
  name = "${var.cluster_name}-etcd-instance-role"
  role = "${aws_iam_role.etcd_role.name}"
}
