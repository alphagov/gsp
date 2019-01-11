data "aws_route53_zone" "zone" {
  name = "${var.cluster_zone}."
}

