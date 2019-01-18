output "ignition-file-ids" {
  value = [
    "${data.ignition_file.kubeconfig.id}",
    "${data.ignition_file.kubeconfig-kubelet.id}",
    "${data.ignition_file.bootstrap-apiserver.id}",
    "${data.ignition_file.bootstrap-controller-manager.id}",
    "${data.ignition_file.bootstrap-scheduler.id}",
    "${data.ignition_file.coredns-cluster-role-binding.id}",
    "${data.ignition_file.coredns-cluster-role.id}",
    "${data.ignition_file.coredns-config-yaml.id}",
    "${data.ignition_file.coredns-deployment.id}",
    "${data.ignition_file.coredns-service-account.id}",
    "${data.ignition_file.coredns-service.id}",
    "${data.ignition_file.aws-iam-authenticator-cfg.id}",
    "${data.ignition_file.aws-iam-authenticator-daemonset.id}",
    "${data.ignition_file.aws-iam-authenticator-kubeconfig.id}",
    "${data.ignition_file.aws-iam-authenticator-secret.id}",
    "${data.ignition_file.default-storage-class.id}",
    "${data.ignition_file.kube-apiserver-secret.id}",
    "${data.ignition_file.kube-apiserver.id}",
    "${data.ignition_file.kube-controller-manager-disruption.id}",
    "${data.ignition_file.kube-controller-manager-role-binding.id}",
    "${data.ignition_file.kube-controller-manager-secret.id}",
    "${data.ignition_file.kube-controller-manager-service-account.id}",
    "${data.ignition_file.kube-controller-manager.id}",
    "${data.ignition_file.kube-proxy-role-binding.id}",
    "${data.ignition_file.kube-proxy-sa.id}",
    "${data.ignition_file.kube-proxy.id}",
    "${data.ignition_file.kube-scheduler-disruption.id}",
    "${data.ignition_file.kube-scheduler.id}",
    "${data.ignition_file.kube-system-rbac-role-binding.id}",
    "${data.ignition_file.kubeconfig-in-cluster.id}",
    "${data.ignition_file.kubernetes-dashboard-role-binding.id}",
    "${data.ignition_file.kubernetes-dashboard-role.id}",
    "${data.ignition_file.kubernetes-dashboard-sa.id}",
    "${data.ignition_file.kubernetes-dashboard-secret.id}",
    "${data.ignition_file.kubernetes-dashboard-svc.id}",
    "${data.ignition_file.kubernetes-dashboard.id}",
    "${data.ignition_file.pod-checkpointer-role-binding.id}",
    "${data.ignition_file.pod-checkpointer-role.id}",
    "${data.ignition_file.pod-checkpointer-sa.id}",
    "${data.ignition_file.pod-checkpointer.id}",
    "${data.ignition_file.tiller-role-binding.id}",
    "${data.ignition_file.tiller-sa.id}",
    "${data.ignition_file.tiller-svc.id}",
    "${data.ignition_file.tiller.id}",
    "${data.ignition_file.calico-bgpconfigurations-crd.id}",
    "${data.ignition_file.calico-bgppeers-crd.id}",
    "${data.ignition_file.calico-clusterinformations-crd.id}",
    "${data.ignition_file.calico-config.id}",
    "${data.ignition_file.calico-felixconfigurations-crd.id}",
    "${data.ignition_file.calico-globalnetworkpolicies-crd.id}",
    "${data.ignition_file.calico-globalnetworksets-crd.id}",
    "${data.ignition_file.calico-hostendpoints-crd.id}",
    "${data.ignition_file.calico-ippools-crd.id}",
    "${data.ignition_file.calico-networkpolicies-crd.id}",
    "${data.ignition_file.calico-node-cluster-role-binding.id}",
    "${data.ignition_file.calico-node-cluster-role.id}",
    "${data.ignition_file.calico-node-service-account.id}",
    "${data.ignition_file.calico.id}",
    "${data.ignition_file.ca-key.id}",
    "${data.ignition_file.ca-crt.id}",
    "${data.ignition_file.apiserver-key.id}",
    "${data.ignition_file.apiserver-crt.id}",
    "${data.ignition_file.service-account-key.id}",
    "${data.ignition_file.service-account-pub.id}",
    "${data.ignition_file.admin-key.id}",
    "${data.ignition_file.admin-crt.id}",
  ]
}

output "admin-kubeconfig" {
  value = "${data.template_file.kubeconfig-user.rendered}"
}

output "kubelet-kubeconfig" {
  value     = "${data.template_file.kubeconfig-kubelet.rendered}"
  sensitive = true
}

output "kube-ca-crt" {
  value = "${tls_self_signed_cert.kube-ca.cert_pem}"
}

output "bootkube_systemd_unit_id" {
  value = "${data.ignition_systemd_unit.bootkube-service.id}"
}
