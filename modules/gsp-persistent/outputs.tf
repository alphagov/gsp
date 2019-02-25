output "private_key_pem" {
  description = "Sealed secrets private key"
  value       = "${tls_private_key.sealed-secrets-key.private_key_pem}"
}

output "cert_pem" {
  description = "Sealed secrets certificate"
  value       = "${tls_self_signed_cert.sealed-secrets-certificate.cert_pem}"
}
