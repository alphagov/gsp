data "aws_caller_identity" "current" {
}

resource "random_password" "concourse_password" {
  length  = 64
  special = false
}

data "template_file" "values" {
  template = file("${path.module}/data/values.yaml")

  vars = {
    cluster_name                     = var.cluster_name
    cluster_domain                   = var.cluster_domain
    cluster_domain_id                = var.cluster_domain_id
    cluster_oidc_provider_url        = jsonencode(module.k8s-cluster.oidc_provider_url)
    cluster_oidc_provider_arn        = jsonencode(module.k8s-cluster.oidc_provider_arn)
    egress_ip_addresses              = jsonencode(var.egress_ips)
    account_name                     = var.account_name
    account_id                       = data.aws_caller_identity.current.account_id
    admin_role_arns                  = jsonencode(concat(var.admin_role_arns, [module.k8s-cluster.aws_node_lifecycle_hook_role_arn]))
    bootstrap_role_arns              = jsonencode(module.k8s-cluster.bootstrap_role_arns)
    concourse_admin_password         = random_password.concourse_password.result
    concourse_teams                  = jsonencode(concat(["main"], var.concourse_teams))
    concourse_main_team_github_teams = jsonencode(var.concourse_main_team_github_teams)
    concourse_worker_count           = var.ci_worker_count
    github_client_id                 = jsonencode(var.github_client_id)
    github_client_secret             = jsonencode(var.github_client_secret)
    github_ca_cert                   = jsonencode(var.github_ca_cert)
    sealed_secrets_public_cert       = base64encode(tls_self_signed_cert.sealed-secrets-certificate.cert_pem)
    sealed_secrets_private_key       = base64encode(tls_private_key.sealed-secrets-key.private_key_pem)
    cloudwatch_log_shipping_role     = aws_iam_role.cloudwatch_log_shipping_role.arn
    service_operator_boundary_arn    = aws_iam_policy.service-operator-managed-role-permissions-boundary.arn
    service_operator_role_arn        = aws_iam_role.gsp-service-operator.arn
    rds_from_worker_security_group   = aws_security_group.rds-from-worker.id
    private_db_subnet_group          = aws_db_subnet_group.private.id
    external_dns_map                 = yamlencode(local.external_dns)
    grafana_default_admin_password   = jsonencode(random_password.grafana_default_admin_password.result)
    eks_version                      = var.eks_version
    cert_manager_role_arn            = aws_iam_role.cert_manager.arn
  }
}
