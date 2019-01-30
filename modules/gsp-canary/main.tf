resource "aws_codecommit_repository" "canary" {
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

  enabled        = 1
  namespace      = "gsp-canary"
  chart_git      = "${aws_codecommit_repository.canary.clone_url_http}"
  chart_ref      = "master"
  chart_path     = "charts/gsp-canary"
  cluster_name   = "${var.cluster_name}"
  cluster_domain = "${var.cluster_name}.${var.dns_zone}"
  addons_dir     = "${var.addons_dir}"
  permitted_roles_regex = "^${aws_iam_role.canary_role.name}$"

  values = <<EOF
    annotations:
      iam.amazonaws.com/role: "${aws_iam_role.canary_role.name}"
    updater:
      helmChartRepoUrl: ${aws_codecommit_repository.canary.clone_url_http}
EOF
}
