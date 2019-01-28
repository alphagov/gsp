data "local_file" "index" {
  filename = "${path.module}/index.js"
}

data "local_file" "logger" {
  filename = "${path.module}/lib/mysplunklogger.js"
}

data "archive_file" "lambda_log_forwarder" {
  type        = "zip"
  output_path = "${path.module}/.terraform/archive_files/lambda_log_forwarder.zip"

  source {
    content  = "${data.local_file.index.content}"
    filename = "index.js"
  }

  source {
    content  = "${data.local_file.logger.content}"
    filename = "lib/mysplunklogger.js"
  }
}

resource "aws_lambda_function" "lambda_log_forwarder" {
  count            = "${var.enabled == 0 ? 0 : 1}"
  filename         = "${data.archive_file.lambda_log_forwarder.output_path}"
  source_code_hash = "${data.archive_file.lambda_log_forwarder.output_base64sha256}"
  function_name    = "${var.cluster_name}_log_forwarder"
  role             = "${aws_iam_role.lambda_log_forwarder.arn}"
  handler          = "index.handler"
  runtime          = "nodejs6.10"
  timeout          = "10"
  memory_size      = "128"
  description      = "A function to forward logs from AWS to a Splunk HEC"

  environment {
    variables = {
      SPLUNK_HEC_TOKEN = "${var.splunk_hec_token}"
      SPLUNK_HEC_URL   = "${var.splunk_hec_url}"
    }
  }
}

resource "aws_cloudwatch_log_group" "lambda_log_forwarder" {
  count             = "${var.enabled == 0 ? 0 : 1}"
  name              = "/aws/lambda/${aws_lambda_function.lambda_log_forwarder.function_name}"
  retention_in_days = 7
}

resource "aws_lambda_permission" "cloudwatch_splunk_logs" {
  count         = "${var.enabled == 0 ? 0 : 1}"
  statement_id  = "${var.cluster_name}_cloudwatch_splunk_logs"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.lambda_log_forwarder.arn}"
  principal     = "logs.eu-west-2.amazonaws.com"
  source_arn    = "${var.cloudwatch_log_group_arn}"
}

resource "aws_cloudwatch_log_subscription_filter" "cloudwatch_splunk_logs" {
  count           = "${var.enabled == 0 ? 0 : 1}"
  depends_on      = ["aws_lambda_permission.cloudwatch_splunk_logs"]
  name            = "${var.cluster_name}_cloudwatch_splunk_logs_subscription_filter"
  destination_arn = "${aws_lambda_function.lambda_log_forwarder.arn}"
  filter_pattern  = ""
  log_group_name  = "${var.cloudwatch_log_group_name}"
}
