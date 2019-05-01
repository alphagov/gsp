output "vpc_id" {
  value = "${aws_vpc.network.id}"
}

output "private_subnet_ids" {
  value = [
    "${module.subnet-0.private_subnet_id}",
    "${module.subnet-1.private_subnet_id}",
    "${module.subnet-2.private_subnet_id}",
  ]
}

output "public_subnet_ids" {
  value = [
    "${module.subnet-0.public_subnet_id}",
    "${module.subnet-1.public_subnet_id}",
    "${module.subnet-2.public_subnet_id}",
  ]
}

output "nat_gateway_public_ips" {
  value = [
    "${module.subnet-0.nat_gateway_public_ip}",
    "${module.subnet-1.nat_gateway_public_ip}",
    "${module.subnet-2.nat_gateway_public_ip}",
  ]
}

output "host_cidr" {
  description = "CIDR IPv4 range to assign to EC2 nodes"
  value       = "${var.host_cidr}"
}
