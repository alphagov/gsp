output "cluster-name" {
  value = "${var.cluster_name}"
}

output "cluster-domain-suffix" {
  value = "${var.cluster_domain}"
}

output "kubeconfig" {
  value = "${module.k8s-cluster.kubeconfig}"
}

output "values" {
  sensitive = true
  value     = "${data.template_file.values.rendered}"
}

/* output "ci-system-release-name" { */
/*   value = "${module.ci-system.release-name}" */
/* } */


/* output "notary-ci-private-key" { */
/*   value = "${module.ci-system.notary-ci-private-key}" */
/* } */


/* output "notary-root-private-key" { */
/*   value = "${module.ci-system.notary-root-private-key}" */
/* } */


/* output "notary-delegation-passphrase" { */
/*   value = "${module.ci-system.notary-delegation-passphrase}" */
/* } */


/* output "notary-root-passphrase" { */
/*   value = "${module.ci-system.notary-root-passphrase}" */
/* } */


/* output "notary-snapshot-passphrase" { */
/*   value = "${module.ci-system.notary-snapshot-passphrase}" */
/* } */


/* output "notary-targets-passphrase" { */
/*   value = "${module.ci-system.notary-targets-passphrase}" */
/* } */


/* output "harbor-password" { */
/*   value = "${module.ci-system.harbor-password}" */
/* } */


/* output "canary-role-arn" { */
/*   value = "${aws_iam_role.canary_role.arn}" */
/* } */


/* output "canary-role-name" { */
/*   value = "${aws_iam_role.canary_role.name}" */
/* } */


/* output "code-commit-repository-arn" { */
/*   value = "${aws_codecommit_repository.canary.arn}" */
/* } */

