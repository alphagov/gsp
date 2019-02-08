output "release-name" {
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
