resource "aws_route53_zone" "subdomain" {
  count         = "${length(var.extra_zones)}"
  name          = "${element(var.extra_zones, count.index)}"
  force_destroy = true
}
