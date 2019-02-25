output "bootstrap-subnet-id" {
  value = "${element(aws_subnet.cluster-private.*.id, 0)}"
}

output "private_subnet_ids" {
  value = ["${aws_subnet.cluster-private.*.id}"]
}

output "public_subnet_ids" {
  value = ["${aws_subnet.cluster-public.*.id}"]
}

output "network_id" {
  value = "${aws_vpc.network.id}"
}

output "nat_gateway_public_ips" {
  value = ["${aws_nat_gateway.cluster.*.public_ip}"]
}

output "cluster-name" {
  value = "${var.cluster_name}"
}

output "host_cidr" {
  description = "CIDR IPv4 range to assign to EC2 nodes"
  value       = "${var.host_cidr}"
}
