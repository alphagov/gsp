# resource "aws_iam_role" "sre" {
#   name = "${var.cluster_name}-sre"
#   assume_role_policy = "${data.aws_iam_policy_document.grant-iam-sre-policy.json}"
# }
# data "aws_iam_policy_document" "grant-iam-sre-policy" {
#   statement {
#     effect  = "Allow"
#     actions = ["sts:AssumeRole"]
#     principals = {
#       type = "AWS"
#       identifiers = ["${var.sre_role_arns}"]
#     }
#     condition {
#       test     = "Bool"
#       variable = "aws:MultiFactorAuthPresent"
#       values   = ["true"]
#     }
#     condition {
#       test     = "IpAddress"
#       variable = "aws:SourceIp"
#       values   = ["${var.gds_external_cidrs}"]
#     }
#   }
# }
# resource "aws_iam_policy_attachment" "cloudwatch-readonly" {
#   name       = "${var.cluster_name}-cloudwatch-readonly-attachment"
#   roles      = ["${aws_iam_role.sre.name}"]
#   policy_arn = "${aws_iam_policy.cloudwatch-readonly.arn}"
# }
# resource "aws_iam_policy" "cloudwatch-readonly" {
#   name = "${var.cluster_name}-cloudwatch-readonly"
#   policy = "${data.aws_iam_policy_document.cloudwatch-readonly.json}"
# }
# data "aws_iam_policy_document" "cloudwatch-readonly" {
#   statement {
#     effect = "Allow"
#     actions = [
#       "autoscaling:Describe*",
#       "cloudwatch:Describe*",
#       "cloudwatch:Get*",
#       "cloudwatch:List*",
#       "logs:Get*",
#       "logs:Describe*",
#       "sns:Get*",
#       "sns:List*",
#     ]
#     resources = ["*"]
#   }
# }
