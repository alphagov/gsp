output "ignition-systemd-unit-ids" {
  value = ["${data.ignition_systemd_unit.amazon-ssm-agent-service.id}"]
}
