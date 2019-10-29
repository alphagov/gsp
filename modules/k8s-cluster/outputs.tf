output "kubeconfig" {
  value = data.template_file.kubeconfig.rendered
}

output "kiam-server-node-instance-role-arn" {
  value = aws_cloudformation_stack.kiam-server-nodes.outputs["NodeInstanceRole"]
}

output "kiam-server-node-instance-role-name" {
  value = replace(data.aws_arn.kiam-server-nodes-role.resource, "role/", "")
}

output "bootstrap_role_arns" {
  value = [
    aws_cloudformation_stack.worker-nodes.outputs["NodeInstanceRole"],
    aws_cloudformation_stack.kiam-server-nodes.outputs["NodeInstanceRole"],
    aws_cloudformation_stack.ci-nodes.outputs["NodeInstanceRole"],
  ]
}

output "worker_tcp_target_group_arn" {
  value = aws_cloudformation_stack.worker-nodes.outputs["TCPTargetGroup"]
}

output "eks-log-group-arn" {
  value = aws_cloudwatch_log_group.eks.arn
}

output "eks-log-group-name" {
  value = aws_cloudwatch_log_group.eks.name
}

output "worker_security_group_id" {
  value = aws_security_group.worker.id
}

output "oidc_provider_url" {
  value = aws_iam_openid_connect_provider.eks.url
}

output "oidc_provider_arn" {
  value = aws_iam_openid_connect_provider.eks.arn
}

