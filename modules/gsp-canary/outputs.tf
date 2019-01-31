output "canary_role_arn" {
  value = "${aws_iam_role.canary_role.arn}"
}

output "canary_role_name" {
  value = "${aws_iam_role.canary_role.name}"
}
