module "lambda_splunk_forwarder" {
  source                    = "../lambda_splunk_forwarder"
  enabled                   = "${var.splunk_enabled}"
  name                      = "pods"
  cloudwatch_log_group_arn  = "${aws_cloudwatch_log_group.logs.arn}"
  cloudwatch_log_group_name = "${aws_cloudwatch_log_group.logs.name}"
  cluster_name              = "${var.cluster_name}"
  splunk_hec_token          = "${var.splunk_hec_token}"
  splunk_hec_url            = "${var.splunk_hec_url}"
  splunk_index              = "${var.splunk_index}"
}
