output "github-deployment-private-key" {
  value = "${tls_private_key.github_deployment_key.private_key_pem}"
}

output "github-deployment-public-key" {
  value = "${tls_private_key.github_deployment_key.public_key_openssh}"
}
