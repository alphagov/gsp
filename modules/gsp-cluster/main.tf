module "etcd-cluster" {
  source = "../etcd-cluster"

  cluster_name            = "${var.cluster_name}"
  dns_zone                = "${var.dns_zone}"
  subnet_ids              = "${aws_subnet.cluster-private.*.id}"
  vpc_id                  = "${aws_vpc.network.id}"
  dns_zone_id             = "${data.aws_route53_zone.zone.zone_id}"
  node_count              = "${var.etcd_node_count}"
  user_data_bucket_name   = "${var.user_data_bucket_name}"
  user_data_bucket_region = "${var.user_data_bucket_region}"
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
}

module "k8s-cluster" {
  source                       = "../k8s-cluster"
  cluster_domain_suffix        = "cluster.local"
  kubelet_kubeconfig           = "${module.bootkube-assets.kubelet-kubeconfig}"
  kube_ca_crt                  = "${module.bootkube-assets.kube-ca-crt}"
  user_data_bucket_name        = "${var.user_data_bucket_name}"
  user_data_bucket_region      = "${var.user_data_bucket_region}"
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
}

module "ingress-system" {
  enabled = "${var.addons["ingress"]}"
  source = "../flux-release"

  namespace      = "ingress-system"
  chart_git      = "https://github.com/alphagov/gsp-ingress-system.git"
  chart_ref      = "master"
  cluster_name   = "${var.cluster_name}"
  cluster_domain = "${var.cluster_name}.${var.dns_zone}"
  addons_dir     = "addons/${var.cluster_name}"
}

module "monitoring-system" {
  source = "../flux-release"

  enabled = "${var.addons["monitoring"]}"
  namespace      = "monitoring-system"
  chart_git      = "https://github.com/alphagov/gsp-monitoring-system.git"
  chart_ref      = "master"
  cluster_name   = "${var.cluster_name}"
  cluster_domain = "${var.cluster_name}.${var.dns_zone}"
  addons_dir     = "addons/${var.cluster_name}"
}

module "secrets-system" {
  source = "../flux-release"

  enabled = "${var.addons["secrets"]}"
  namespace      = "secrets-system"
  chart_git      = "https://github.com/alphagov/gsp-secrets-system.git"
  chart_ref      = "master"
  cluster_name   = "${var.cluster_name}"
  cluster_domain = "${var.cluster_name}.${var.dns_zone}"
  addons_dir     = "addons/${var.cluster_name}"
}

module "ci-system" {
  source = "..//flux-release"

  enabled = "${var.addons["ci"]}"
  namespace      = "ci-system"
  chart_git      = "https://github.com/alphagov/gsp-ci-system.git"
  chart_ref      = "master"
  cluster_name   = "${var.cluster_name}"
  cluster_domain = "${var.cluster_name}.${var.dns_zone}"
  addons_dir     = "addons/${var.cluster_name}"
}

resource "aws_codecommit_repository" "canary" {
  count = "${var.addons["canary"] ? 1 : 0}"
  repository_name = "canary.${var.cluster_name}.${var.dns_zone}"

  provisioner "local-exec" {
    command = "${path.module}/scripts/initialise_canary_helm_codecommit.sh"

    environment {
      SOURCE_REPO_URL     = "https://github.com/alphagov/gsp-canary-chart"
      CODECOMMIT_REPO_URL = "${aws_codecommit_repository.canary.clone_url_http}"
    }
  }
}

module "canary-system" {
  source = "../flux-release"

  enabled = "${var.addons["canary"]}"
  namespace  = "gsp-canary"
  chart_git  = "${var.addons["canary"] ? element(concat(aws_codecommit_repository.canary.*.clone_url_http, list("")), 0) : ""}"
  chart_ref  = "master"
  chart_path = "charts/gsp-canary"
  cluster_name = ""
  cluster_domain = "${var.cluster_name}.${var.dns_zone}"
  addons_dir     = "addons/${var.cluster_name}"
  values = <<EOF
    updater:
    helmChartRepoUrl: ${var.addons["canary"] ? element(concat(aws_codecommit_repository.canary.*.clone_url_http, list("")), 0) : ""}
EOF
}
