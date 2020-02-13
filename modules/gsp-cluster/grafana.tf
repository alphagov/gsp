resource "random_password" "grafana_default_admin_password" {
  length  = 40
  special = false
}

