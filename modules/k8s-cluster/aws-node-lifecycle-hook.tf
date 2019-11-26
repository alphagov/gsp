resource "aws_iam_role" "aws-node-lifecycle-hook" {
  name               = "${var.cluster_name}_aws-node-lifecycle-hook"
  assume_role_policy = data.aws_iam_policy_document.aws-node-lifecycle-hook_assume_role_policy.json
}

resource "aws_iam_policy_attachment" "aws-node-lifecycle-hook" {
  name       = "${var.cluster_name}_aws-node-lifecycle-hook_attachment"
  roles      = [aws_iam_role.aws-node-lifecycle-hook.name]
  policy_arn = aws_iam_policy.aws-node-lifecycle-hook.arn
}

resource "aws_iam_policy" "aws-node-lifecycle-hook" {
  name        = "${var.cluster_name}-aws-node-lifecycle-hook"
  description = "Policy for aws-node-lifecycle-hook function"
  policy      = data.aws_iam_policy_document.aws-node-lifecycle-hook.json
}

data "aws_iam_policy_document" "aws-node-lifecycle-hook" {
  statement {
    effect = "Allow"

    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = ["*"]
  }

  statement {
    effect = "Allow"

    actions = [
      "eks:DescribeCluster"
    ]

    resources = [
      aws_eks_cluster.eks-cluster.arn
    ]
  }

  statement {
    effect = "Allow"

    actions = [
      "autoscaling:CompleteLifecycleAction",
      "autoscaling:RecordLifecycleActionHeartbeat"
    ]

    resources = ["*"]
  }
}

data "aws_iam_policy_document" "aws-node-lifecycle-hook_assume_role_policy" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_lambda_function" "aws-node-lifecycle-hook" {
  filename         = "${path.module}/aws-node-lifecycle-hook.zip"
  source_code_hash = filebase64sha256("${path.module}/aws-node-lifecycle-hook.zip")
  function_name    = "${var.cluster_name}-aws-node-lifecycle-hook"

  role        = aws_iam_role.aws-node-lifecycle-hook.arn
  handler     = "aws-node-lifecycle-hook"
  runtime     = "go1.x"
  timeout     = "600"
  memory_size = "128"
  description = "A function to handle EKS worker node lifecycle changes"

  environment {
    variables = {
      CLUSTER_NAME = var.cluster_name
    }
  }
}

resource "aws_lambda_permission" "aws-node-lifecycle-hook" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.aws-node-lifecycle-hook.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.aws-node-lifecycle-hook.arn
}

resource "aws_cloudwatch_event_rule" "aws-node-lifecycle-hook" {
  name        = "${var.cluster_name}-asg-lifecycle-hook"
  description = "Execute ASG lifecycle logic"

  event_pattern = <<PATTERN
{
  "detail-type": [
    "EC2 Instance-terminate Lifecycle Action"
  ],
  "source": [
    "aws.autoscaling"
  ],
  "detail": {
    "AutoScalingGroupName": ${jsonencode(
  concat(
    [
      for stack in aws_cloudformation_stack.worker-nodes-per-az : lookup(stack.outputs, "AutoScalingGroupName", "")
    ],
    [
      lookup(aws_cloudformation_stack.kiam-server-nodes.outputs, "AutoScalingGroupName", ""),
      lookup(aws_cloudformation_stack.ci-nodes.outputs, "AutoScalingGroupName", "")
    ]
  )
)}
  }
}
PATTERN
}

resource "aws_cloudwatch_event_target" "aws-node-lifecycle-hook" {
  rule      = aws_cloudwatch_event_rule.aws-node-lifecycle-hook.name
  target_id = "ASGLifecycleLambda"
  arn       = aws_lambda_function.aws-node-lifecycle-hook.arn
}

