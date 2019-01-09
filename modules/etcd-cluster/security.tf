resource "aws_security_group" "etcd" {
  name        = "${var.cluster_name}-etcd"
  description = "${var.cluster_name} etcd security group"

  vpc_id = "${var.vpc_id}"

  tags = "${map("Name", "${var.cluster_name}-controller")}"
}

resource "aws_security_group_rule" "ssh" {
  security_group_id = "${aws_security_group.etcd.id}"

  type        = "ingress"
  protocol    = "tcp"
  from_port   = 22
  to_port     = 22
  cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "etcd" {
  security_group_id = "${aws_security_group.etcd.id}"

  type        = "ingress"
  protocol    = "tcp"
  from_port   = 2379
  to_port     = 2380
  cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "egress" {
  security_group_id = "${aws_security_group.etcd.id}"

  type        = "egress"
  protocol    = "-1"
  from_port   = 0
  to_port     = 0
  cidr_blocks = ["0.0.0.0/0"]
}
