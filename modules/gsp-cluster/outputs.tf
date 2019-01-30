output "bootstrap-base-userdata-source" {
  value = "s3://${var.user_data_bucket_name}${aws_s3_bucket_object.bootstrap-user-data.id}"
}

output "bootstrap-base-userdata-verification" {
  value = "sha512-${sha512(data.ignition_config.bootstrap.rendered)}"
}

output "controller-security-group-ids" {
  value = ["${module.k8s-cluster.controller-security-group-ids}"]
}

output "bootstrap-subnet-id" {
  value = "${element(aws_subnet.cluster-private.*.id, 0)}"
}

output "controller-instance-profile-name" {
  value = "${module.k8s-cluster.controller-instance-profile-name}"
}

output "apiserver-lb-target-group-arn" {
  value = "${aws_lb_target_group.controllers.arn}"
}

output "dns-service-ip" {
  value = "${module.k8s-cluster.dns-service-ip}"
}

output "cluster-name" {
  value = "${var.cluster_name}"
}

output "k8s_tag" {
  value = "${var.k8s_tag}"
}

output "user_data_bucket_name" {
  value = "${var.user_data_bucket_name}"
}

output "user_data_bucket_region" {
  value = "${var.user_data_bucket_region}"
}

output "cluster-domain-suffix" {
  value = "${var.cluster_name}.${var.dns_zone}"
}

output "kubelet-kubeconfig" {
  value     = "${module.bootkube-assets.kubelet-kubeconfig}"
  sensitive = true
}

output "admin-kubeconfig" {
  value = "${module.bootkube-assets.admin-kubeconfig}"
}

output "kube-ca-crt" {
  value = "${module.bootkube-assets.kube-ca-crt}"
}

output "ci-system-release-name" {
  value = "${module.ci-system.release-name}"
}

output "notary-ci-private-key" {
  value = "${base64encode(tls_private_key.notary_ci_key.private_key_pem)}"
}

output "notary-root-private-key" {
  value = "${base64encode(tls_private_key.notary_root_key.private_key_pem)}"
}

output "notary-delegation-passphrase" {
  value = "${base64encode(random_string.notary_passphrase_delegation.result)}"
}

output "notary-root-passphrase" {
  value = "${base64encode(random_string.notary_passphrase_root.result)}"
}

output "notary-snapshot-passphrase" {
  value = "${base64encode(random_string.notary_passphrase_snapshot.result)}"
}

output "notary-targets-passphrase" {
  value = "${base64encode(random_string.notary_passphrase_targets.result)}"
}

output "harbor-password" {
  value = "${base64encode(random_string.harbor_password.result)}"
}

output "github-deployment-public-key" {
  value = "${tls_private_key.github_deployment_key.public_key_openssh}"
}
