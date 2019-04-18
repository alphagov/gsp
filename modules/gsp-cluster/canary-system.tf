resource "aws_iam_role" "canary_role" {
  name               = "${var.cluster_name}-canary"
  description        = "Role the gsp-canary process assumes"
  assume_role_policy = "${data.aws_iam_policy_document.assume_canary_role.json}"
}

resource "aws_iam_policy" "canary_code_commit" {
  name        = "${var.cluster_name}-canary-code-commit"
  description = "Policy for the gsp-canary code commit"
  policy      = "${data.aws_iam_policy_document.canary_code_commit.json}"
}

data "aws_iam_policy_document" "assume_canary_role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals = {
      type        = "AWS"
      identifiers = ["${aws_iam_role.kiam_server_role.arn}"]
    }
  }
}

data "aws_iam_policy_document" "canary_code_commit" {
  statement {
    actions = [
      "codecommit:GitPull",
      "codecommit:GitPush",
    ]

    resources = ["${aws_codecommit_repository.canary.arn}"]
  }
}

resource "aws_iam_policy_attachment" "canary_code_commit" {
  name       = "${var.cluster_name}-canary-code-commit"
  roles      = ["${aws_iam_role.canary_role.name}"]
  policy_arn = "${aws_iam_policy.canary_code_commit.arn}"
}

resource "aws_codecommit_repository" "canary" {
  repository_name = "canary.${var.cluster_name}.${var.account_name}"

  provisioner "local-exec" {
    command = "${path.module}/scripts/initialise_canary_helm_codecommit.sh"

    environment {
      SOURCE_REPO_URL          = "https://github.com/alphagov/gsp-canary-chart"
      CODECOMMIT_REPO_URL      = "${aws_codecommit_repository.canary.clone_url_http}"
      CODECOMMIT_INIT_ROLE_ARN = "${var.codecommit_init_role_arn}"
    }
  }
}

/* module "canary-system" { */
/*   source = "../flux-release" */
/*   namespace             = "gsp-canary" */
/*   chart_git             = "${aws_codecommit_repository.canary.clone_url_http}" */
/*   chart_ref             = "master" */
/*   chart_path            = "charts/gsp-canary" */
/*   permitted_roles_regex = "^${aws_iam_role.canary_role.name}$" */
/*   values = <<EOF */
/*     annotations: */
/*       iam.amazonaws.com/role: "${aws_iam_role.canary_role.name}" */
/*     updater: */
/*       helmChartRepoUrl: ${aws_codecommit_repository.canary.clone_url_http} */
/* EOF */
/* } */
/* module "gsp-canary" { */
/*   addons_dir               = "addons/${var.cluster_name}" */
/*   canary_role_assumer_arn  = "${aws_iam_role.kiam_server_role.arn}" */
/*   codecommit_init_role_arn = "${var.codecommit_init_role_arn}" */
/* } */

