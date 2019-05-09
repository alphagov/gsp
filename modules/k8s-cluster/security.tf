data "aws_vpc" "private" {
  id = "${var.vpc_id}"
}

resource "aws_security_group" "controller" {
  name        = "${var.cluster_name}-controller"
  description = "${var.cluster_name} controller security group"

  vpc_id = "${var.vpc_id}"

  tags = "${map("Name", "${var.cluster_name}-controller")}"
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

resource "aws_security_group_rule" "worker-nodes-from-vpc" {
  security_group_id = "${aws_cloudformation_stack.worker-nodes.outputs["NodeSecurityGroup"]}"

  type      = "ingress"
  protocol  = "-1"
  from_port = 0
  to_port   = 0

  cidr_blocks = ["${data.aws_vpc.private.cidr_block}"]
}

resource "aws_security_group_rule" "worker-nodes-nodeport-30080" {
  security_group_id = "${aws_cloudformation_stack.worker-nodes.outputs["NodeSecurityGroup"]}"

  type      = "ingress"
  protocol  = "tcp"
  from_port = 30080
  to_port   = 30080

  cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "worker-nodes-nodeport-30443" {
  security_group_id = "${aws_cloudformation_stack.worker-nodes.outputs["NodeSecurityGroup"]}"

  type      = "ingress"
  protocol  = "tcp"
  from_port = 30443
  to_port   = 30443

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
