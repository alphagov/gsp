resource "aws_flow_log" "vpc_flow_log" {
  iam_role_arn         = "${aws_iam_role.cloudwatch_vpc_flow_log_shipping.arn}"
  log_destination      = "${aws_cloudwatch_log_group.vpc_flow_log.arn}"
  log_destination_type = "cloud-watch-logs"
  traffic_type         = "ALL"
  vpc_id               = "${var.vpc_id}"
}

resource "aws_cloudwatch_log_group" "vpc_flow_log" {
  name              = "${var.cluster_domain}_vpc_flow_log"
  retention_in_days = 30
}

resource "aws_iam_role" "cloudwatch_vpc_flow_log_shipping" {
  name = "${var.cluster_name}_cloudwatch_vpc_flow_log_shipping"

  assume_role_policy = "${data.aws_iam_policy_document.cloudwatch_vpc_flow_log_assume_role.json}"
}

resource "aws_iam_policy" "cloudwatch_vpc_flow_log_shipping" {
  name        = "${var.cluster_name}_cloudwatch_vpc_flow_log_shipping"
  description = "Send logs to Clouwatch"

  policy = "${data.aws_iam_policy_document.cloudwatch_vpc_flow_log.json}"
}

resource "aws_iam_policy_attachment" "cloudwatch_vpc_flow_log_shipping" {
  name       = "${var.cluster_name}_cloudwatch_vpc_flow_log_shipping"
  roles      = ["${aws_iam_role.cloudwatch_vpc_flow_log_shipping.name}"]
  policy_arn = "${aws_iam_policy.cloudwatch_vpc_flow_log_shipping.arn}"
}

data "aws_iam_policy_document" "cloudwatch_vpc_flow_log_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["vpc-flow-log.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "cloudwatch_vpc_flow_log" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
      "logs:DescribeLogGroups",
      "logs:DescribeLogStreams",
    ]

    resources = [
      "*",
    ]
  }
}
