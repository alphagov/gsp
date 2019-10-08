data "aws_caller_identity" "current" {}

data "template_file" "values" {
  template = "${file("${path.module}/data/values.yaml")}"

  vars {
    cluster_name                     = "${var.cluster_name}"
    cluster_domain                   = "${var.cluster_domain}"
    cluster_domain_id                = "${var.cluster_domain_id}"
    cluster_oidc_provider_url        = "${jsonencode(module.k8s-cluster.oidc_provider_url)}"
    cluster_oidc_provider_arn        = "${jsonencode(module.k8s-cluster.oidc_provider_arn)}"
    egress_ip_addresses              = "${jsonencode(var.egress_ips)}"
    account_name                     = "${var.account_name}"
    account_id                       = "${data.aws_caller_identity.current.account_id}"
    admin_role_arns                  = "${jsonencode(var.admin_role_arns)}"
    admin_user_arns                  = "${jsonencode(var.admin_user_arns)}"
    sre_role_arns                    = "${jsonencode(var.sre_role_arns)}"
    sre_user_arns                    = "${jsonencode(var.sre_user_arns)}"
    bootstrap_role_arns              = "${jsonencode(module.k8s-cluster.bootstrap_role_arns)}"
    cluster_autoscaler_role_name     = "${aws_iam_role.cluster_autoscaler.name}"
    concourse_admin_password         = "${random_string.concourse_password.result}"
    concourse_teams                  = "${jsonencode(concat(list("main"), var.concourse_teams))}"
    concourse_main_team_github_teams = "${jsonencode(var.concourse_main_team_github_teams)}"
    concourse_worker_count           = "${var.ci_worker_count}"
    concourse_iam_role_name          = "${aws_iam_role.concourse.name}"
    github_client_id                 = "${jsonencode(var.github_client_id)}"
    github_client_secret             = "${jsonencode(var.github_client_secret)}"
    github_ca_cert                   = "${jsonencode(var.github_ca_cert)}"
    grafana_iam_role_name            = "${aws_iam_role.grafana.name}"
    harbor_admin_password            = "${jsonencode(random_string.harbor_password.result)}"
    harbor_secret_key                = "${jsonencode(random_string.harbor_secret_key.result)}"
    harbor_bucket_id                 = "${aws_s3_bucket.ci-system-harbor-registry-storage.id}"
    harbor_bucket_region             = "${aws_s3_bucket.ci-system-harbor-registry-storage.region}"
    harbor_iam_role_name             = "${jsonencode(aws_iam_role.harbor.name)}"
    notary_root_key                  = "${jsonencode(tls_private_key.notary_root_key.private_key_pem)}"
    notary_ca_pem                    = "${jsonencode(tls_self_signed_cert.notary_root_ca.cert_pem)}"
    notary_cert_pem                  = "${jsonencode(tls_locally_signed_cert.notary_cert.cert_pem)}"
    notary_delegation_key            = "${jsonencode(tls_private_key.notary_ci_key.private_key_pem)}"
    notary_root_passphrase           = "${jsonencode(random_string.notary_passphrase_root.result)}"
    notary_targets_passphrase        = "${jsonencode(random_string.notary_passphrase_targets.result)}"
    notary_snapshot_passphrase       = "${jsonencode(random_string.notary_passphrase_snapshot.result)}"
    notary_delegation_passphrase     = "${jsonencode(random_string.notary_passphrase_delegation.result)}"
    sealed_secrets_public_cert       = "${base64encode(tls_self_signed_cert.sealed-secrets-certificate.cert_pem)}"
    sealed_secrets_private_key       = "${base64encode(tls_private_key.sealed-secrets-key.private_key_pem)}"
    kiam_server_role_arn             = "${aws_iam_role.kiam_server_role.arn}"
    kiam_ca_cert_b64e_pem            = "${base64encode(tls_self_signed_cert.kiam_ca.cert_pem)}"
    kiam_server_cert_b64e_pem        = "${base64encode(tls_locally_signed_cert.kiam_server.cert_pem)}"
    kiam_server_key_b64e_pem         = "${base64encode(tls_private_key.kiam_server.private_key_pem)}"
    kiam_agent_cert_b64e_pem         = "${base64encode(tls_locally_signed_cert.kiam_agent.cert_pem)}"
    kiam_agent_key_b64e_pem          = "${base64encode(tls_private_key.kiam_agent.private_key_pem)}"
    cloudwatch_log_shipping_role     = "${aws_iam_role.cloudwatch_log_shipping_role.name}"
    cloudwatch_log_group_name        = "${aws_cloudwatch_log_group.logs.name}"
    service_operator_boundary_arn    = "${aws_iam_policy.service-operator-managed-role-permissions-boundary.arn}"
    rds_from_worker_security_group   = "${aws_security_group.rds-from-worker.id}"
    private_db_subnet_group          = "${aws_db_subnet_group.private.id}"
    external_dns_iam_role_name       = "${aws_iam_role.external_dns.name}"
    grafana_default_admin_password   = "${random_password.grafana_default_admin_password.result}"
    eks_version                      = "${var.eks_version}"

    permitted_roles_regex = "^(${join("|", list(
      aws_iam_role.cloudwatch_log_shipping_role.name,
      aws_iam_role.cluster_autoscaler.name,
      aws_iam_role.concourse.name,
      aws_iam_role.grafana.name,
      aws_iam_role.gsp-service-operator.name,
      aws_iam_role.harbor.name,
      aws_iam_role.external_dns.name,
    ))})$"

    istio_system_permitted_roles_regex = "^(${join("|", list(
      aws_iam_role.cert_manager.name,
    ))})$"
  }
}
