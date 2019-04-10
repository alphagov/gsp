module "k8s-cluster" {
  source               = "../k8s-cluster"
  vpc_id               = "${var.vpc_id}"
  subnet_ids           = ["${concat(var.private_subnet_ids, var.public_subnet_ids)}"]
  cluster_name         = "${var.cluster_name}"
  worker_count         = "${var.worker_count}"
  worker_instance_type = "${var.worker_instance_type}"
  ci_worker_count         = "${var.ci_worker_count}"
  ci_worker_instance_type = "${var.ci_worker_instance_type}"
  admin_role_arns      = "${var.admin_role_arns}"
  sre_role_arns        = "${var.sre_user_arns}"

  apiserver_allowed_cidrs = ["${concat(
      formatlist("%s/32", var.nat_gateway_public_ips),
      var.gds_external_cidrs,
  )}"]
}

resource "local_file" "aws-auth" {
  filename = "addons/${var.cluster_name}/aws-auth.yaml"
  content  = "${module.k8s-cluster.aws-auth}"
}

locals {
  default_addons = {
    ingress    = 1
    monitoring = 1
    secrets    = 1
    ci         = 1
    splunk     = 0
  }

  enabled_addons = "${merge(local.default_addons, var.addons)}"
}

data "template_file" "flux-reporter" {
  template = "${file("${path.module}/data/flux-reporter.yaml")}"

  vars {
    namespace      = "flux-system"
    cluster_domain = "${var.cluster_name}.${var.dns_zone}"
  }
}

resource "local_file" "flux-reporter" {
  filename = "addons/${var.cluster_name}/flux-reporter.yaml"
  content  = "${data.template_file.flux-reporter.rendered}"
}

resource "local_file" "cert-manager-crds" {
  count    = "${local.enabled_addons["ingress"]}"
  filename = "addons/${var.cluster_name}/cert-manager-crds.yaml"
  content  = "${file("${path.module}/data/cert-manager-crds.yaml")}"
}

module "monitoring-system" {
  source = "../flux-release"

  enabled               = "${local.enabled_addons["monitoring"]}"
  namespace             = "monitoring-system"
  chart_git             = "https://github.com/alphagov/gsp-monitoring-system.git"
  chart_ref             = "master"
  cluster_name          = "${var.cluster_name}"
  cluster_domain        = "${var.cluster_name}.${var.dns_zone}"
  addons_dir            = "addons/${var.cluster_name}"
  permitted_roles_regex = "^${aws_iam_role.cloudwatch_log_shipping_role.name}$"

  values = <<EOF
    fluentd-cloudwatch:
      logGroupName: "${var.cluster_name}.${var.dns_zone}"
      awsRole: "${aws_iam_role.cloudwatch_log_shipping_role.name}"
    prometheus-operator:
      prometheus:
        prometheusSpec:
          externalLabels:
            clustername: "${var.cluster_name}.${var.dns_zone}"
            product: "${var.account_name}"
            deployment: gsp
EOF
}

resource "aws_cloudwatch_log_group" "logs" {
  count             = "${local.enabled_addons["monitoring"] ? 1 : 0}"
  name              = "${var.cluster_name}.${var.dns_zone}"
  retention_in_days = 30
}

module "lambda_splunk_forwarder" {
  source = "../lambda_splunk_forwarder"

  enabled                   = "${local.enabled_addons["splunk"]}"
  name                      = "pods"
  cloudwatch_log_group_arn  = "${aws_cloudwatch_log_group.logs.arn}"
  cloudwatch_log_group_name = "${aws_cloudwatch_log_group.logs.name}"
  cluster_name              = "${var.cluster_name}"
  splunk_hec_token          = "${var.splunk_hec_token}"
  splunk_hec_url            = "${var.splunk_hec_url}"
  splunk_index              = "${var.splunk_index}"
}

module "secrets-system" {
  source = "../flux-release"

  enabled        = "${local.enabled_addons["secrets"]}"
  namespace      = "secrets-system"
  chart_git      = "https://github.com/alphagov/gsp-secrets-system.git"
  chart_ref      = "master"
  cluster_name   = "${var.cluster_name}"
  cluster_domain = "${var.cluster_name}.${var.dns_zone}"
  addons_dir     = "addons/${var.cluster_name}"

  values = <<EOF
    encryption:
      public_certificate: ${base64encode(var.sealed_secrets_cert_pem)}
      private_key: ${base64encode(var.sealed_secrets_private_key_pem)}
EOF
}

module "kiam-system" {
  source = "../flux-release"

  enabled        = 1
  namespace      = "kiam-system"
  chart_git      = "https://github.com/alphagov/gsp-kiam-system"
  chart_ref      = "master"
  cluster_name   = "${var.cluster_name}"
  cluster_domain = "${var.cluster_name}.${var.dns_zone}"
  addons_dir     = "addons/${var.cluster_name}"

  values = <<EOF
    kiam:
      server:
        assumeRoleArn: "${aws_iam_role.kiam_server_role.arn}"
      agent:
        host:
          interface: "eni+"
EOF
}

module "group-role-bindings" {
  source = "../group-role-bindings"

  namespaces = ["${var.dev_namespaces}"]
  addons_dir = "addons/${var.cluster_name}"
  group_name = "dev"
}

module "gsp-canary" {
  source                   = "../gsp-canary"
  cluster_name             = "${var.cluster_name}"
  dns_zone                 = "${var.dns_zone}"
  addons_dir               = "addons/${var.cluster_name}"
  canary_role_assumer_arn  = "${aws_iam_role.kiam_server_role.arn}"
  codecommit_init_role_arn = "${var.codecommit_init_role_arn}"
}

module "ci-system" {
  source                 = "../gsp-ci"
  enabled                = "${local.enabled_addons["ci"]}"
  cluster_name           = "${var.cluster_name}"
  dns_zone               = "${var.dns_zone}"
  harbor_role_asumer_arn = "${aws_iam_role.kiam_server_role.arn}"

  github_teams         = "${var.github_teams}"
  github_client_id     = "${var.github_client_id}"
  github_client_secret = "${var.github_client_secret}"
  github_ca_cert       = "${var.github_ca_cert}"
  concourse_teams      = "${var.concourse_teams}"
}

resource "local_file" "role" {
  filename = "addons/${var.cluster_name}/sre-cluster-role.yaml"
  content  = "${file("${path.module}/data/sre-cluster-role.yaml")}"
}

resource "local_file" "aws-ssm-agent-daemonset" {
  filename = "addons/${var.cluster_name}/aws-ssm-agent-daemonset.yaml"
  content  = "${file("${path.module}/data/aws-ssm-agent-daemonset.yaml")}"
}
