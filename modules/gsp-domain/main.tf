provider "aws" {}

provider "aws" {
  alias  = "apex"
  region = "eu-west-2"
}

data "aws_route53_zone" "apex" {
  provider = "aws.apex"
  name     = "${var.existing_zone}"
}

resource "aws_route53_zone" "subdomain" {
  name          = "${var.delegated_zone}"
  force_destroy = true
}

resource "aws_route53_record" "ns" {
  provider = "aws.apex"
  zone_id  = "${data.aws_route53_zone.apex.zone_id}"
  name     = "${var.delegated_zone}"
  type     = "NS"
  ttl      = "30"

  records = [
    "${aws_route53_zone.subdomain.name_servers.0}",
    "${aws_route53_zone.subdomain.name_servers.1}",
    "${aws_route53_zone.subdomain.name_servers.2}",
    "${aws_route53_zone.subdomain.name_servers.3}",
  ]
}
