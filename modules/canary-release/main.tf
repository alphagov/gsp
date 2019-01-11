resource "aws_codecommit_repository" "canary" {
  repository_name = "gsp-canary-chart-${var.cluster_id}"

  provisioner "local-exec" {
    command = "../../../scripts/initialise_canary_helm_codecommit.sh"

    environment {
      SOURCE_REPO_URL     = "https://github.com/alphagov/gsp-canary-chart"
      CODECOMMIT_REPO_URL = "${aws_codecommit_repository.canary.clone_url_http}"
    }
  }
}

module "canary-release" {
  source = "../flux-release"

  namespace  = "gsp-canary"
  chart_git  = "${aws_codecommit_repository.canary.clone_url_http}"
  chart_ref  = "master"
  chart_path = "charts/gsp-canary"
  cluster_name = ""
  cluster_domain = "${var.cluster_id}"
  values = <<EOF
    updater:
      helmChartRepoUrl: ${aws_codecommit_repository.canary.clone_url_http}
EOF
}
