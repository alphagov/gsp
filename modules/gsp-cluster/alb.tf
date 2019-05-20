resource "aws_security_group" "ingress" {
  name        = "${var.cluster_name}-ingress"
  description = "${var.cluster_name} ingress (ALB) security group"

  vpc_id = "${var.vpc_id}"

  tags = "${map(
    "Name", "${var.cluster_name}-controller",
  )}"
}

resource "aws_security_group_rule" "ingress-https" {
  security_group_id = "${aws_security_group.ingress.id}"
  type              = "ingress"
  protocol          = "tcp"
  from_port         = 443
  to_port           = 443
  cidr_blocks       = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "ingress-http" {
  security_group_id = "${aws_security_group.ingress.id}"
  type              = "ingress"
  protocol          = "tcp"
  from_port         = 80
  to_port           = 80
  cidr_blocks       = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "ingress-outbound" {
  security_group_id = "${aws_security_group.ingress.id}"
  type              = "egress"
  protocol          = "-1"
  from_port         = 0
  to_port           = 0
  cidr_blocks       = ["0.0.0.0/0"]
}

resource "aws_lb" "ingress" {
  name               = "${var.cluster_name}-ingress"
  load_balancer_type = "application"
  security_groups    = ["${aws_security_group.ingress.id}"]

  subnet_mapping {
    subnet_id = "${var.public_subnet_ids[0]}"
  }

  subnet_mapping {
    subnet_id = "${var.public_subnet_ids[1]}"
  }

  subnet_mapping {
    subnet_id = "${var.public_subnet_ids[2]}"
  }

  tags = "${map("Name", "${var.cluster_name}-ingress")}"
}

resource "aws_acm_certificate" "default" {
  domain_name       = "*.${var.cluster_domain}"
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_route53_record" "default" {
  name    = "${aws_acm_certificate.default.domain_validation_options.0.resource_record_name}"
  type    = "${aws_acm_certificate.default.domain_validation_options.0.resource_record_type}"
  zone_id = "${var.cluster_domain_id}"
  records = ["${aws_acm_certificate.default.domain_validation_options.0.resource_record_value}"]
  ttl     = 60
}

resource "aws_acm_certificate_validation" "default" {
  certificate_arn         = "${aws_acm_certificate.default.arn}"
  validation_record_fqdns = ["${aws_route53_record.default.fqdn}"]
}

resource "aws_lb_listener" "ingress-https" {
  load_balancer_arn = "${aws_lb.ingress.arn}"
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = "${aws_acm_certificate_validation.default.certificate_arn}"

  default_action {
    type             = "forward"
    target_group_arn = "${module.k8s-cluster.worker_http_target_group_arn}"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_lb_listener" "ingress-http" {
  load_balancer_arn = "${aws_lb.ingress.arn}"
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type = "redirect"

    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

resource "aws_route53_record" "ingress-root" {
  zone_id = "${var.cluster_domain_id}"
  name    = "${var.cluster_domain}."
  type    = "A"

  alias {
    name                   = "${aws_lb.ingress.dns_name}"
    zone_id                = "${aws_lb.ingress.zone_id}"
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "ingress-wildcard" {
  zone_id = "${var.cluster_domain_id}"
  name    = "*.${var.cluster_domain}."
  type    = "A"

  alias {
    name                   = "${aws_lb.ingress.dns_name}"
    zone_id                = "${aws_lb.ingress.zone_id}"
    evaluate_target_health = true
  }
}
