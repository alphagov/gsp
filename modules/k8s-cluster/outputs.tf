output "controller-security-group-ids" {
  value = ["${aws_security_group.controller.id}"]
}

output "worker-security-group-ids" {
  value = ["${aws_security_group.worker.id}"]
}

output "controller-instance-profile-name" {
  value = "${aws_iam_instance_profile.controller_profile.name}"
}

output "controller-role-arn" {
  value = "${aws_iam_role.controller_role.arn}"
}

output "dns-service-ip" {
  value = "${cidrhost(var.service_cidr, 10)}"
}
