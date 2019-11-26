variable "subnet_cidr_map" {
  type = map(string)
}

variable "source_security_group_id" {
  type = string
}

variable "cluster_name" {
  type = string
}

variable "splunk" {
  default = 0
}

variable "splunk_hec_url" {
  type = string
}

variable "splunk_hec_token" {
  type = string
}

variable "splunk_index" {
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

module "lambda_splunk_forwarder" {
  source = "../lambda_splunk_forwarder"

  enabled                   = var.splunk
  name                      = "hsm"
  cloudwatch_log_group_arn  = "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/cloudhsm/${aws_cloudhsm_v2_cluster.cluster.cluster_id}:*"
  cloudwatch_log_group_name = "/aws/cloudhsm/${aws_cloudhsm_v2_cluster.cluster.cluster_id}"
  cluster_name              = var.cluster_name
  splunk_hec_token          = var.splunk_hec_token
  splunk_hec_url            = var.splunk_hec_url
  splunk_index              = var.splunk_index
}
resource "aws_cloudwatch_log_subscription_filter" "hsm_logs" {
  count           = var.cls_destination_enabled == "1" ? 1 : 0
  name            = "hsm_logs"
  log_group_name  = "/aws/cloudhsm/${aws_cloudhsm_v2_cluster.cluster.cluster_id}"
  filter_pattern  = ""
  destination_arn = var.cls_destination_arn
}
