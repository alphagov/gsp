variable "subnet_cidr_map" {
  type = map(string)
}

variable "source_security_group_id" {
  type = string
}

variable "cluster_name" {
  type = string
}

variable "cls_destination_enabled" {
  default = "0"
  type    = string
}

variable "cls_destination_arn" {
  default = ""
  type    = string
}

data "aws_caller_identity" "current" {
}

data "aws_region" "current" {
}

resource "aws_cloudhsm_v2_cluster" "cluster" {
  hsm_type = "hsm1.medium"

  subnet_ids = keys(var.subnet_cidr_map)

  tags = {
    Name = "${var.cluster_name}-hsm-cluster"
  }
}

resource "aws_security_group_rule" "hsm-worker-ingress" {
  security_group_id        = aws_cloudhsm_v2_cluster.cluster.security_group_id
  type                     = "ingress"
  from_port                = 2223
  to_port                  = 2225
  protocol                 = "tcp"
  source_security_group_id = var.source_security_group_id
}

resource "aws_cloudwatch_log_subscription_filter" "hsm_logs" {
  count           = var.cls_destination_enabled == "1" ? 1 : 0
  name            = "hsm_logs"
  log_group_name  = "/aws/cloudhsm/${aws_cloudhsm_v2_cluster.cluster.cluster_id}"
  filter_pattern  = ""
  destination_arn = var.cls_destination_arn
}
