module "k8s_lambda_splunk_forwarder" {
  source                    = "../lambda_splunk_forwarder"
  enabled                   = "${var.splunk_enabled}"
  name                      = "pods"
  cloudwatch_log_group_arn  = "${aws_cloudwatch_log_group.logs.arn}"
  cloudwatch_log_group_name = "${aws_cloudwatch_log_group.logs.name}"
  cluster_name              = "${var.cluster_name}"
  splunk_hec_token          = "${var.k8s_splunk_hec_token}"
  splunk_hec_url            = "${var.splunk_hec_url}"
  splunk_index              = "${var.k8s_splunk_index}"
}

module "eks_lambda_splunk_forwarder" {
  source                    = "../lambda_splunk_forwarder"
  enabled                   = "${var.splunk_enabled}"
  name                      = "eks"
  cloudwatch_log_group_arn  = "${module.k8s-cluster.eks-log-group-arn}"
  cloudwatch_log_group_name = "${module.k8s-cluster.eks-log-group-name}"
  cluster_name              = "${var.cluster_name}"
  splunk_hec_token          = "${var.k8s_splunk_hec_token}"
  splunk_hec_url            = "${var.splunk_hec_url}"
  splunk_index              = "${var.k8s_splunk_index}"
}

module "vpc_flow_log_lambda_splunk_forwarder" {
  source                    = "../lambda_splunk_forwarder"
  enabled                   = "${var.splunk_enabled}"
  name                      = "vpc_flow_log"
  cloudwatch_log_group_arn  = "${aws_cloudwatch_log_group.vpc_flow_log.arn}"
  cloudwatch_log_group_name = "${aws_cloudwatch_log_group.vpc_flow_log.name}"
  cluster_name              = "${var.cluster_name}"
  splunk_hec_token          = "${var.vpc_flow_log_splunk_hec_token}"
  splunk_hec_url            = "${var.splunk_hec_url}"
  splunk_index              = "${var.vpc_flow_log_splunk_index}"
}
