data "aws_caller_identity" "current" {}

data "template_file" "values" {
  template = "${file("${path.module}/data/values.yaml")}"

  vars {
    cluster_name                        = "${var.cluster_name}"
    cluster_domain                      = "${var.cluster_domain}"
    cluster_domain_id                   = "${var.cluster_domain_id}"
    account_name                        = "${var.account_name}"
    account_id                          = "${data.aws_caller_identity.current.account_id}"
    admin_role_arns                     = "${jsonencode(var.admin_role_arns)}"
    admin_user_arns                     = "${jsonencode(var.admin_user_arns)}"
    sre_role_arns                       = "${jsonencode(var.sre_role_arns)}"
    sre_user_arns                       = "${jsonencode(var.sre_user_arns)}"
    bootstrap_role_arns                 = "${jsonencode(module.k8s-cluster.bootstrap_role_arns)}"
    concourse_admin_password            = "${random_string.concourse_password.result}"
    concourse_teams                     = "${jsonencode(concat(list("main"), var.concourse_teams))}"
    concourse_main_team_github_teams    = "${jsonencode(var.concourse_main_team_github_teams)}"
    concourse_worker_count              = "${var.ci_worker_count}"
    concourse_iam_role_name             = "${aws_iam_role.concourse.name}"
    github_client_id                    = "${jsonencode(var.github_client_id)}"
    github_client_secret                = "${jsonencode(var.github_client_secret)}"
    github_ca_cert                      = "${jsonencode(var.github_ca_cert)}"
    harbor_admin_password               = "${jsonencode(random_string.harbor_password.result)}"
    harbor_secret_key                   = "${jsonencode(random_string.harbor_secret_key.result)}"
    harbor_bucket_id                    = "${aws_s3_bucket.ci-system-harbor-registry-storage.id}"
    harbor_bucket_region                = "${aws_s3_bucket.ci-system-harbor-registry-storage.region}"
    harbor_iam_role_name                = "${jsonencode(aws_iam_role.harbor.name)}"
    notary_root_key                     = "${jsonencode(tls_private_key.notary_root_key.private_key_pem)}"
    notary_ca_pem                       = "${jsonencode(tls_self_signed_cert.notary_root_ca.cert_pem)}"
    notary_cert_pem                     = "${jsonencode(tls_locally_signed_cert.notary_cert.cert_pem)}"
    notary_delegation_key               = "${jsonencode(tls_private_key.notary_ci_key.private_key_pem)}"
    notary_root_passphrase              = "${jsonencode(random_string.notary_passphrase_root.result)}"
    notary_targets_passphrase           = "${jsonencode(random_string.notary_passphrase_targets.result)}"
    notary_snapshot_passphrase          = "${jsonencode(random_string.notary_passphrase_snapshot.result)}"
    notary_delegation_passphrase        = "${jsonencode(random_string.notary_passphrase_delegation.result)}"
    sealed_secrets_public_cert          = "${base64encode(tls_self_signed_cert.sealed-secrets-certificate.cert_pem)}"
    sealed_secrets_private_key          = "${base64encode(tls_private_key.sealed-secrets-key.private_key_pem)}"
    kiam_server_role_arn                = "${aws_iam_role.kiam_server_role.arn}"
    kiam_restart_after_deploy_hack_uuid = "${uuid()}"
    cloudwatch_log_shipping_role        = "${aws_iam_role.cloudwatch_log_shipping_role.name}"
    cloudwatch_log_group_name           = "${aws_cloudwatch_log_group.logs.name}"

    permitted_roles_regex = "^(${join("|", list(
      aws_iam_role.harbor.name,
      aws_iam_role.concourse.name,
      aws_iam_role.cloudwatch_log_shipping_role.name,
    ))})$"
  }
}
