data "ignition_systemd_unit" "amazon-ssm-agent-service" {
  name = "amazon-ssm-agent.service"

  content = "${file("${path.module}/data/amazon-ssm-agent.service")}"
}
