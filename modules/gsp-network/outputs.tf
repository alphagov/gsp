output "bootstrap-subnet-id" {
  value = "${element(aws_subnet.cluster-private.*.id, 0)}"
}

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

output "cluster-name" {
  value = "${var.cluster_name}"
}

output "host_cidr" {
  description = "CIDR IPv4 range to assign to EC2 nodes"
  value       = "${var.host_cidr}"
}

// Workaround for https://github.com/hashicorp/terraform/issues/12570
output "private-subnet-count" {
  value = "${length(data.aws_availability_zones.all.names)}"
}
