output "zone_id" {
  value = "${aws_route53_zone.subdomain.zone_id}"
}

output "name" {
  value = "${aws_route53_zone.subdomain.name}"
}
