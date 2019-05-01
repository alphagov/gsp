output "nat_gateway_public_ip" {
  value = "${aws_eip.egress.public_ip}"
}

output "private_subnet_id" {
  value = "${aws_subnet.private.id}"
}

output "public_subnet_id" {
  value = "${aws_subnet.public.id}"
}
