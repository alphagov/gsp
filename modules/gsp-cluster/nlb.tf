data "aws_eip" "ingress-0" {
  public_ip = "${var.ingress_ips[0]}"
}

data "aws_eip" "ingress-1" {
  public_ip = "${var.ingress_ips[1]}"
}

data "aws_eip" "ingress-2" {
  public_ip = "${var.ingress_ips[2]}"
}

resource "aws_lb" "ingress" {
  name                             = "${var.cluster_name}-ingress"
  load_balancer_type               = "network"
  enable_cross_zone_load_balancing = true

  subnet_mapping {
    subnet_id     = "${var.public_subnet_ids[0]}"
    allocation_id = "${data.aws_eip.ingress-0.id}"
  }

  subnet_mapping {
    subnet_id     = "${var.public_subnet_ids[1]}"
    allocation_id = "${data.aws_eip.ingress-1.id}"
  }

  subnet_mapping {
    subnet_id     = "${var.public_subnet_ids[2]}"
    allocation_id = "${data.aws_eip.ingress-2.id}"
  }

  tags = "${map("Name", "${var.cluster_name}-ingress")}"
}

resource "aws_lb_listener" "ingress-https" {
  load_balancer_arn = "${aws_lb.ingress.arn}"
  port              = "443"
  protocol          = "TCP"

  default_action {
    type             = "forward"
    target_group_arn = "${module.k8s-cluster.worker_https_target_group_arn}"
  }
}

resource "aws_lb_listener" "ingress-http" {
  load_balancer_arn = "${aws_lb.ingress.arn}"
  port              = "80"
  protocol          = "TCP"

  default_action {
    type             = "forward"
    target_group_arn = "${module.k8s-cluster.worker_http_target_group_arn}"
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
