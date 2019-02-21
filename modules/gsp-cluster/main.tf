module "etcd-cluster" {
  source = "../etcd-cluster"

  cluster_name            = "${var.cluster_name}"
  dns_zone                = "${var.dns_zone}"
  subnet_ids              = "${aws_subnet.cluster-private.*.id}"
  vpc_id                  = "${aws_vpc.network.id}"
  dns_zone_id             = "${data.aws_route53_zone.zone.zone_id}"
  node_count              = "${var.etcd_node_count}"
  user_data_bucket_name   = "${var.user_data_bucket_name}"
  instance_type           = "${var.etcd_instance_type}"
  s3_user_data_policy_arn = "${aws_iam_policy.s3-user-data-policy.arn}"
}

module "bootkube-assets" {
  source                      = "../bootkube-ignition"
  apiserver_address           = "${aws_route53_record.apiserver.fqdn}"
  cluster_domain_suffix       = "cluster.local"
  etcd_servers                = ["${module.etcd-cluster.etcd_servers}"]
  k8s_tag                     = "${var.k8s_tag}"
  cluster_name                = "${var.cluster_name}"
  cluster_id                  = "${var.cluster_name}.${var.dns_zone}"
  etcd_ca_cert_pem            = "${module.etcd-cluster.ca_cert_pem}"
  etcd_client_private_key_pem = "${module.etcd-cluster.client_private_key_pem}"
  etcd_client_cert_pem        = "${module.etcd-cluster.client_cert_pem}"
  admin_role_arns             = ["${var.admin_role_arns}"]
  dev_role_arns               = ["${aws_iam_role.dev.arn}"]
}

module "k8s-cluster" {
  source                       = "../k8s-cluster"
  cluster_domain_suffix        = "cluster.local"
  kubelet_kubeconfig           = "${module.bootkube-assets.kubelet-kubeconfig}"
  kube_ca_crt                  = "${module.bootkube-assets.kube-ca-crt}"
  user_data_bucket_name        = "${var.user_data_bucket_name}"
  vpc_id                       = "${aws_vpc.network.id}"
  subnet_ids                   = ["${aws_subnet.cluster-private.*.id}"]
  controller_target_group_arns = ["${aws_lb_target_group.controllers.arn}"]

  worker_target_group_arns = [
    "${aws_lb_target_group.workers-http.arn}",
    "${aws_lb_target_group.workers-https.arn}",
  ]

  cluster_name             = "${var.cluster_name}"
  k8s_tag                  = "${var.k8s_tag}"
  controller_count         = "${var.controller_count}"
  worker_count             = "${var.worker_count}"
  controller_instance_type = "${var.controller_instance_type}"
  worker_instance_type     = "${var.worker_instance_type}"
  s3_user_data_policy_arn  = "${aws_iam_policy.s3-user-data-policy.arn}"

  apiserver_allowed_cidrs = ["${concat(
      list(aws_vpc.network.cidr_block),
      formatlist("%s/32", aws_nat_gateway.cluster.*.public_ip),
      var.gds_external_cidrs,
  )}"]
}

locals {
  default_addons = {
    ingress    = 1
    monitoring = 1
    secrets    = 1
    ci         = 0
    splunk     = 0
  }

  enabled_addons = "${merge(local.default_addons, var.addons)}"
}

data "template_file" "flux" {
  template = "${file("${path.module}/data/flux.yaml")}"

  vars {
    namespace             = "flux-system"
    aws_role_name         = "${module.gsp-canary.canary_role_name}"
    permitted_roles_regex = "^${module.gsp-canary.canary_role_name}$"
  }
}

resource "local_file" "flux" {
  filename = "addons/${var.cluster_name}/flux.yaml"
  content  = "${data.template_file.flux.rendered}"
}

module "ingress-system" {
  enabled = "${local.enabled_addons["ingress"]}"
  source  = "../flux-release"

  namespace      = "ingress-system"
  chart_git      = "https://github.com/alphagov/gsp-ingress-system.git"
  chart_ref      = "master"
  cluster_name   = "${var.cluster_name}"
  cluster_domain = "${var.cluster_name}.${var.dns_zone}"
  addons_dir     = "addons/${var.cluster_name}"

  extra_namespace_configuration = <<EOF
  labels:
    certmanager.k8s.io/disable-validation: "true"
EOF

  values = <<EOF
    webhook:
      enabled: false
EOF
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
      public_certificate: ${base64encode(tls_self_signed_cert.sealed-secrets-certificate.cert_pem)}
      private_key: ${base64encode(tls_private_key.sealed-secrets-key.private_key_pem)}
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
EOF
}

resource "tls_private_key" "github_deployment_key" {
  algorithm = "RSA"
  rsa_bits  = "4096"
}

resource "local_file" "ci-secrets" {
  count    = "${local.enabled_addons["ci"]}"
  filename = "addons/${var.cluster_name}/ci-secrets.yaml"
  content  = "${data.template_file.ci-secrets.rendered}"
}

data "template_file" "ci-secrets" {
  count    = "${local.enabled_addons["ci"]}"
  template = "${file("${path.module}/data/ci-secrets.yaml")}"

  vars = {
    namespace   = "ci-system-main"
    private_key = "${base64encode(tls_private_key.github_deployment_key.private_key_pem)}"
  }
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
}
