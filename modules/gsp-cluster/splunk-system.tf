module "k8s_lambda_splunk_forwarder" {
  source                    = "../lambda_splunk_forwarder"
  enabled                   = var.splunk_enabled
  name                      = "pods"
  cloudwatch_log_group_arn  = aws_cloudwatch_log_group.logs.arn
  cloudwatch_log_group_name = aws_cloudwatch_log_group.logs.name
  cluster_name              = var.cluster_name
  splunk_hec_token          = var.k8s_splunk_hec_token
  splunk_hec_url            = var.splunk_hec_url
  splunk_index              = var.k8s_splunk_index
}
resource "aws_cloudwatch_log_subscription_filter" "legacy_logs" {
  count           = var.cls_destination_enabled == "1" ? 1 : 0
  name            = "legacy_logs"
  log_group_name  = aws_cloudwatch_log_group.logs.name
  filter_pattern  = ""
  destination_arn = var.cls_destination_arn
}

module "k8s_app_lambda_splunk_forwarder" {
  source                    = "../lambda_splunk_forwarder"
  enabled                   = var.splunk_enabled
  name                      = "application"
  cloudwatch_log_group_arn  = aws_cloudwatch_log_group.application_logs.arn
  cloudwatch_log_group_name = aws_cloudwatch_log_group.application_logs.name
  cluster_name              = var.cluster_name
  splunk_hec_token          = var.k8s_splunk_hec_token
  splunk_hec_url            = var.splunk_hec_url
  splunk_index              = var.k8s_splunk_index
}
resource "aws_cloudwatch_log_subscription_filter" "application_logs" {
  count           = var.cls_destination_enabled == "1" ? 1 : 0
  name            = "application_logs"
  log_group_name  = aws_cloudwatch_log_group.application_logs.name
  filter_pattern  = ""
  destination_arn = var.cls_destination_arn
}

module "k8s_host_lambda_splunk_forwarder" {
  source                    = "../lambda_splunk_forwarder"
  enabled                   = var.splunk_enabled
  name                      = "host"
  cloudwatch_log_group_arn  = aws_cloudwatch_log_group.host_logs.arn
  cloudwatch_log_group_name = aws_cloudwatch_log_group.host_logs.name
  cluster_name              = var.cluster_name
  splunk_hec_token          = var.k8s_splunk_hec_token
  splunk_hec_url            = var.splunk_hec_url
  splunk_index              = var.k8s_splunk_index
}
resource "aws_cloudwatch_log_subscription_filter" "host_logs" {
  count           = var.cls_destination_enabled == "1" ? 1 : 0
  name            = "host_logs"
  log_group_name  = aws_cloudwatch_log_group.host_logs.name
  filter_pattern  = ""
  destination_arn = var.cls_destination_arn
}

module "k8s_dataplane_lambda_splunk_forwarder" {
  source                    = "../lambda_splunk_forwarder"
  enabled                   = var.splunk_enabled
  name                      = "dataplane"
  cloudwatch_log_group_arn  = aws_cloudwatch_log_group.dataplane_logs.arn
  cloudwatch_log_group_name = aws_cloudwatch_log_group.dataplane_logs.name
  cluster_name              = var.cluster_name
  splunk_hec_token          = var.k8s_splunk_hec_token
  splunk_hec_url            = var.splunk_hec_url
  splunk_index              = var.k8s_splunk_index
}
resource "aws_cloudwatch_log_subscription_filter" "dataplane_logs" {
  count           = var.cls_destination_enabled == "1" ? 1 : 0
  name            = "dataplane_logs"
  log_group_name  = aws_cloudwatch_log_group.dataplane_logs.name
  filter_pattern  = ""
  destination_arn = var.cls_destination_arn
}

module "eks_lambda_splunk_forwarder" {
  source                    = "../lambda_splunk_forwarder"
  enabled                   = var.splunk_enabled
  name                      = "eks"
  cloudwatch_log_group_arn  = module.k8s-cluster.eks-log-group-arn
  cloudwatch_log_group_name = module.k8s-cluster.eks-log-group-name
  cluster_name              = var.cluster_name
  splunk_hec_token          = var.k8s_splunk_hec_token
  splunk_hec_url            = var.splunk_hec_url
  splunk_index              = var.k8s_splunk_index
}
resource "aws_cloudwatch_log_subscription_filter" "eks_logs" {
  count           = var.cls_destination_enabled == "1" ? 1 : 0
  name            = "eks_logs"
  log_group_name  = module.k8s-cluster.eks-log-group-name
  filter_pattern  = ""
  destination_arn = var.cls_destination_arn
}
