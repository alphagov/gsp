output "ignition_file_ids" {
  value = [
    "${data.ignition_file.kubelet-kubeconfig.id}",
    "${data.ignition_file.kube-ca-crt.id}",
  ]
}

output "ignition_systemd_unit_ids" {
  value = ["${data.ignition_systemd_unit.kubelet-service.id}"]
}
