output "controller-security-group-ids" {
  value = ["${aws_security_group.controller.id}"]
}

output "controller-instance-profile-name" {
  value = "${aws_iam_instance_profile.controller_profile.name}"
}

output "dns-service-ip" {
  value = "${cidrhost(var.service_cidr, 10)}"
}
