resource "aws_lambda_function" "lambda_log_forwarder" {
  count            = var.enabled == "0" ? 0 : 1
  filename         = "${path.module}/cyber-cloudwatch-fluentd-to-hec.zip"
  source_code_hash = filebase64sha256("${path.module}/cyber-cloudwatch-fluentd-to-hec.zip")
  function_name    = "${var.cluster_name}_${var.name}_log_forwarder"

  role        = "${aws_iam_role.lambda_log_forwarder[0].arn}"
  handler     = "lambda_function.lambda_handler"
  runtime     = "python3.6"
  timeout     = "120"
  memory_size = "128"
  description = "A function to forward logs from AWS to a Splunk HEC using a manual zip of https://github.com/alphagov/cyber-cloudwatch-fluentd-to-hec"

  environment {
    variables = {
      SPLUNK_HEC_TOKEN = var.splunk_hec_token
      SPLUNK_HEC_URL   = var.splunk_hec_url
      SPLUNK_INDEX     = var.splunk_index
    }
  }
}

resource "aws_cloudwatch_log_group" "lambda_log_forwarder" {
  count = var.enabled == "0" ? 0 : 1

  name              = "/aws/lambda/${aws_lambda_function.lambda_log_forwarder[0].function_name}"
  retention_in_days = 7
}

resource "aws_lambda_permission" "cloudwatch_splunk_logs" {
  count        = var.enabled == "0" ? 0 : 1
  statement_id = "${var.cluster_name}_cloudwatch_splunk_logs"
  action       = "lambda:InvokeFunction"

  function_name = "${aws_lambda_function.lambda_log_forwarder[0].arn}"
  principal     = "logs.eu-west-2.amazonaws.com"
  source_arn    = var.cloudwatch_log_group_arn
}

resource "aws_cloudwatch_log_subscription_filter" "cloudwatch_splunk_logs" {
  count      = var.enabled == "0" ? 0 : 1
  depends_on = [aws_lambda_permission.cloudwatch_splunk_logs]
  name       = "${var.cluster_name}_${var.name}_cloudwatch_splunk_logs_subscription_filter"

  destination_arn = "${aws_lambda_function.lambda_log_forwarder[0].arn}"
  filter_pattern  = ""
  log_group_name  = var.cloudwatch_log_group_name
}

