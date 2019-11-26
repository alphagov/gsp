resource "aws_cloudwatch_log_subscription_filter" "legacy_logs" {
  count           = var.cls_destination_enabled == "1" ? 1 : 0
  name            = "legacy_logs"
  log_group_name  = aws_cloudwatch_log_group.logs.name
  filter_pattern  = ""
  destination_arn = var.cls_destination_arn
}

resource "aws_cloudwatch_log_subscription_filter" "application_logs" {
  count           = var.cls_destination_enabled == "1" ? 1 : 0
  name            = "application_logs"
  log_group_name  = aws_cloudwatch_log_group.application_logs.name
  filter_pattern  = ""
  destination_arn = var.cls_destination_arn
}

resource "aws_cloudwatch_log_subscription_filter" "host_logs" {
  count           = var.cls_destination_enabled == "1" ? 1 : 0
  name            = "host_logs"
  log_group_name  = aws_cloudwatch_log_group.host_logs.name
  filter_pattern  = ""
  destination_arn = var.cls_destination_arn
}

resource "aws_cloudwatch_log_subscription_filter" "dataplane_logs" {
  count           = var.cls_destination_enabled == "1" ? 1 : 0
  name            = "dataplane_logs"
  log_group_name  = aws_cloudwatch_log_group.dataplane_logs.name
  filter_pattern  = ""
  destination_arn = var.cls_destination_arn
}

resource "aws_cloudwatch_log_subscription_filter" "eks_logs" {
  count           = var.cls_destination_enabled == "1" ? 1 : 0
  name            = "eks_logs"
  log_group_name  = module.k8s-cluster.eks-log-group-name
  filter_pattern  = ""
  destination_arn = var.cls_destination_arn
}
