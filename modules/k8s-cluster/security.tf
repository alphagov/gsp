data "aws_vpc" "private" {
  id = "${var.vpc_id}"
}

resource "aws_security_group" "controller" {
  name        = "${var.cluster_name}-controller"
  description = "${var.cluster_name} controller security group"

  vpc_id = "${var.vpc_id}"

  tags = "${map(
    "Name", "${var.cluster_name}-controller",
    "kubernetes.io/cluster/${var.cluster_name}", "owned",
  )}"
}

resource "aws_security_group_rule" "controller-apiserver-cidrs" {
  security_group_id = "${aws_security_group.controller.id}"

  type        = "ingress"
  protocol    = "tcp"
  from_port   = 443
  to_port     = 443
  cidr_blocks = ["${var.apiserver_allowed_cidrs}"]
}

resource "aws_security_group_rule" "controller-egress" {
  security_group_id = "${aws_security_group.controller.id}"

  type        = "egress"
  protocol    = "-1"
  from_port   = 0
  to_port     = 0
  cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "kiam-server-from-vpc" {
  security_group_id = "${aws_cloudformation_stack.kiam-server-nodes.outputs["NodeSecurityGroup"]}"

  type        = "ingress"
  protocol    = "-1"
  from_port   = 0
  to_port     = 0
  cidr_blocks = ["${data.aws_vpc.private.cidr_block}"]
}

resource "aws_security_group_rule" "ci-nodes-from-vpc" {
  security_group_id = "${aws_cloudformation_stack.ci-nodes.outputs["NodeSecurityGroup"]}"

  type        = "ingress"
  protocol    = "-1"
  from_port   = 0
  to_port     = 0
  cidr_blocks = ["${data.aws_vpc.private.cidr_block}"]
}

resource "aws_security_group" "node" {
  name        = "${var.cluster_name}-node"
  description = "${var.cluster_name} node security group.  All nodes should be in this security group."

  vpc_id = "${var.vpc_id}"

  tags = "${map(
    "Name", "${var.cluster_name}-node",
    "kubernetes.io/cluster/${var.cluster_name}", "owned",
  )}"
}

resource "aws_security_group_rule" "node-egress" {
  security_group_id = "${aws_security_group.node.id}"

  type        = "egress"
  protocol    = "-1"
  from_port   = 0
  to_port     = 0
  cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "nodes-from-vpc" {
  security_group_id = "${aws_security_group.node.id}"

  type      = "ingress"
  protocol  = "-1"
  from_port = 0
  to_port   = 0

  cidr_blocks = ["${data.aws_vpc.private.cidr_block}"]
}

resource "aws_security_group_rule" "nodes-from-controller" {
  security_group_id = "${aws_security_group.node.id}"

  type      = "ingress"
  protocol  = "tcp"
  from_port = 1025
  to_port   = 65535

  source_security_group_id = "${aws_security_group.controller.id}"
}

resource "aws_security_group_rule" "controller-to-nodes" {
  security_group_id = "${aws_security_group.controller.id}"

  type      = "egress"
  protocol  = "tcp"
  from_port = 1025
  to_port   = 65535

  source_security_group_id = "${aws_security_group.node.id}"
}

resource "aws_security_group_rule" "controller-from-nodes" {
  security_group_id = "${aws_security_group.controller.id}"

  type      = "ingress"
  protocol  = "tcp"
  from_port = 443
  to_port   = 443

  source_security_group_id = "${aws_security_group.node.id}"
}

resource "aws_security_group" "worker" {
  name        = "${var.cluster_name}-worker"
  description = "${var.cluster_name} worker node security group - ie nodes that tenant pods will run on."

  vpc_id = "${var.vpc_id}"

  tags = "${map(
    "Name", "${var.cluster_name}-worker",
  )}"
}

resource "aws_security_group_rule" "workers-from-public" {
  security_group_id = "${aws_security_group.worker.id}"

  type      = "ingress"
  protocol  = "tcp"
  from_port = 31390
  to_port   = 31390

  cidr_blocks = ["0.0.0.0/0"]
}
