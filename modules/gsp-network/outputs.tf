output "private_subnet_ids" {
  value = ["${aws_subnet.cluster-private.*.id}"]
}

output "public_subnet_ids" {
  value = ["${aws_subnet.cluster-public.*.id}"]
}

output "vpc_id" {
  value = "${aws_vpc.network.id}"
}

output "nat_gateway_public_ips" {
  value = ["${aws_nat_gateway.cluster.*.public_ip}"]
}

output "host_cidr" {
  description = "CIDR IPv4 range to assign to EC2 nodes"
  value       = "${var.host_cidr}"
}
