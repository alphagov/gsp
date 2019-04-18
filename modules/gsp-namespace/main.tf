resource "tls_private_key" "github_deployment_key" {
  algorithm = "RSA"
  rsa_bits  = "4096"
}

/* resource "aws_iam_role" "dev" { */
/*   name = "${var.cluster_name}-${var.name}" */
/*   assume_role_policy = "${data.aws_iam_policy_document.grant-iam-dev.json}" */
/* } */


/* data "aws_iam_policy_document" "grant-iam-dev" { */
/*   statement { */
/*     effect  = "Allow" */
/*     actions = ["sts:AssumeRole"] */
/*     principals = { */
/*       type        = "AWS" */
/*       identifiers = ["${concat(var.admin_role_arns, var.dev_user_arns)}"] */
/*     } */
/*     condition { */
/*       test     = "Bool" */
/*       variable = "aws:MultiFactorAuthPresent" */
/*       values   = ["true"] */
/*     } */
/*     condition { */
/*       test     = "IpAddress" */
/*       variable = "aws:SourceIp" */
/*       values   = ["${var.gds_external_cidrs}"] */
/*     } */
/*   } */
/* } */

