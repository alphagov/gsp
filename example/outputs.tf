output "bootstrap-base-userdata-source" {
    value = "${module.gsp-cluster.bootstrap-base-userdata-source}"
}

output "bootstrap-base-userdata-verification" {
    value = "${module.gsp-cluster.bootstrap-base-userdata-verification}"
}

output "user-data-bucket-name" {
    value = "${var.user_data_bucket_name}"
}

output "user-data-bucket-region" {
    value = "${var.user_data_bucket_region}"
}

output "cluster-name" {
    value = "${var.cluster_name}"
}

output "controller-security-group-ids" {
    value = ["${module.gsp-cluster.controller-security-group-ids}"]
}

output "bootstrap-subnet-id" {
    value = "${module.gsp-cluster.bootstrap-subnet-id}"
}

output "controller-instance-profile-name" {
    value = "${module.gsp-cluster.controller-instance-profile-name}"
}

output "apiserver-lb-target-group-arn" {
    value = "${module.gsp-cluster.apiserver-lb-target-group-arn}"
}

output "dns-service-ip" {
    value = "${module.gsp-cluster.dns-service-ip}"
}

output "cluster-domain-suffix" {
    value = "${module.gsp-cluster.cluster-domain-suffix}"
}

output "k8s-tag" {
    value = "${var.k8s_tag}"
}

output "kubelet-kubeconfig" {
    value = "${module.gsp-cluster.kubelet-kubeconfig}"
}

output "admin-kubeconfig" {
    value = "${module.gsp-cluster.admin-kubeconfig}"
}

output "kube-ca-crt" {
    value = "${module.gsp-cluster.kube-ca-crt}"
}
