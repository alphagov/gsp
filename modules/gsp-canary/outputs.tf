output "canary-role-arn" {
  value = "${aws_iam_role.canary_role.arn}"
}

output "canary-role-name" {
  value = "${aws_iam_role.canary_role.name}"
}

output "code-commit-repository-arn" {
  value = "${aws_codecommit_repository.canary.arn}"
}
