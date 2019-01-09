output "bootkube-ignition-file-ids" {
  value = [
    "${data.ignition_file.etcd-client-ca-crt.id}",
    "${data.ignition_file.etcd-server-ca-crt.id}",
    "${data.ignition_file.etcd-peer-ca-crt.id}",
    "${data.ignition_file.etcd-client-key.id}",
    "${data.ignition_file.etcd-client-crt.id}",
    "${data.ignition_file.etcd-server-key.id}",
    "${data.ignition_file.etcd-server-crt.id}",
    "${data.ignition_file.etcd-peer-key.id}",
    "${data.ignition_file.etcd-peer-crt.id}",
  ]
}

output "etcd_servers" {
  value = "${aws_route53_record.etcds.*.fqdn}"
}

output "ca_cert_pem" {
  value = "${tls_self_signed_cert.etcd-ca.cert_pem}"
}

output "client_cert_pem" {
  value = "${tls_locally_signed_cert.etcd-client.cert_pem}"
}

output "client_private_key_pem" {
  value = "${tls_private_key.etcd-client.private_key_pem}"
}
