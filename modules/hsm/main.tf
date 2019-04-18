variable "subnet_ids" {
  type = "list"
}

variable "subnet_count" {
  type = "string"
}

variable "cluster_name" {
  type = "string"
}

data "aws_subnet" "vpc" {
  count = "${var.subnet_count}"
  id = "${var.subnet_ids[count.index]}"
}

variable "splunk" {
  default = 0
}

variable "splunk_hec_url" {
  type = "string"
}

variable "splunk_hec_token" {
  type = "string"
}

variable "splunk_index" {
  type = "string"
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

resource "aws_cloudhsm_v2_cluster" "cluster" {
  hsm_type   = "hsm1.medium"
  subnet_ids = ["${var.subnet_ids}"]

  tags = {
    Name = "${var.cluster_name}-hsm-cluster"
  }
}

resource "aws_security_group_rule" "hsm-worker-ingress" {
  security_group_id        = "${aws_cloudhsm_v2_cluster.cluster.security_group_id}"
  type                     = "ingress"
  from_port                = 2223
  to_port                  = 2225
  protocol                 = "tcp"
  cidr_blocks = ["${data.aws_subnet.vpc.*.cidr_block}"]
}

# We can only create one HSM in Terraform rather than the multiple we require for high availability as you must create
# a single HSM, initialise and activate it (which is done manually) before you can create more as they are clones of the
# first HSM. The other HSMs will need to be created after the Terraform apply
# Manual steps to initalise and activate the HSM can be followed from
# https://docs.aws.amazon.com/cloudhsm/latest/userguide/configure-sg.html onwards
resource "aws_cloudhsm_v2_hsm" "cloudhsm_v2_hsm" {
  subnet_id  = "${aws_cloudhsm_v2_cluster.cluster.subnet_ids[0]}"
  cluster_id = "${aws_cloudhsm_v2_cluster.cluster.cluster_id}"
}

module "lambda_splunk_forwarder" {
  source = "git::https://github.com/alphagov/gsp-terraform-ignition//modules/lambda_splunk_forwarder?ref=7e322a41cb4c4d45cd36d3a4ee301cd127b95a9a"

  enabled                   = "${var.splunk}"
  name                      = "hsm"
  cloudwatch_log_group_arn  = "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/cloudhsm/${aws_cloudhsm_v2_cluster.cluster.cluster_id}:*"
  cloudwatch_log_group_name = "/aws/cloudhsm/${aws_cloudhsm_v2_cluster.cluster.cluster_id}"
  cluster_name              = "${var.cluster_name}"
  splunk_hec_token          = "${var.splunk_hec_token}"
  splunk_hec_url            = "${var.splunk_hec_url}"
  splunk_index              = "${var.splunk_index}"
}
