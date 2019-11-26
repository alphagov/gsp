output "vpc_id" {
  value = aws_vpc.network.id
}

output "private_subnet_ids" {
  value = [
    module.subnet-0.private_subnet_id,
    module.subnet-1.private_subnet_id,
    module.subnet-2.private_subnet_id,
  ]
}

output "private_subnet_cidr_mapping" {
  value = {
    "${module.subnet-0.private_subnet_id}" = module.subnet-0.private_subnet_cidr
    "${module.subnet-1.private_subnet_id}" = module.subnet-1.private_subnet_cidr
    "${module.subnet-2.private_subnet_id}" = module.subnet-2.private_subnet_cidr
  }
}

output "public_subnet_ids" {
  value = [
    module.subnet-0.public_subnet_id,
    module.subnet-1.public_subnet_id,
    module.subnet-2.public_subnet_id,
  ]
}

output "egress_ips" {
  value = [
    module.subnet-0.egress_ip,
    module.subnet-1.egress_ip,
    module.subnet-2.egress_ip,
  ]
}

output "ingress_ips" {
  value = [
    module.subnet-0.ingress_ip,
    module.subnet-1.ingress_ip,
    module.subnet-2.ingress_ip,
  ]
}

output "cidr_block" {
  description = "CIDR IPv4 range of the VPC"
  value       = aws_vpc.network.cidr_block
}

output "availability_zones" {
  value = [
    module.subnet-0.availability_zone,
    module.subnet-1.availability_zone,
    module.subnet-2.availability_zone,
  ]
}
