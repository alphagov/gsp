output "egress_ip" {
  value = aws_eip.egress.public_ip
}

output "egress_id" {
  value = aws_eip.egress.id
}

output "ingress_ip" {
  value = aws_eip.ingress.public_ip
}

output "ingress_id" {
  value = aws_eip.ingress.id
}

output "private_subnet_id" {
  value = aws_subnet.private.id
}

output "private_subnet_cidr" {
  value = aws_subnet.private.cidr_block
}

output "public_subnet_id" {
  value = aws_subnet.public.id
}

output "availability_zone" {
  value = var.availability_zone
}
