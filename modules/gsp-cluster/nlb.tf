resource "aws_lb" "ingress-nlb" {
  count = "${var.enable_nlb == 1 ? 1 : 0 }"

  name               = "${var.cluster_name}-ingress-nlb"
  load_balancer_type = "network"

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

resource "aws_lb_listener" "ingress-nlb" {
  count = "${var.enable_nlb == 1 ? 1 : 0 }"

  load_balancer_arn = "${aws_lb.ingress-nlb[0].arn}"
  protocol          = "TCP"
  port              = "443"

  default_action {
    type             = "forward"
    target_group_arn = "${module.k8s-cluster.worker_tcp_target_group_arn}"
  }
}

resource "aws_route53_record" "ingress-nlb" {
  count = "${var.enable_nlb == 1 ? 1 : 0 }"

  zone_id = "${var.cluster_domain_id}"
  name    = "nlb.${var.cluster_domain}."
  type    = "A"

  alias {
    name                   = "${aws_lb.ingress-nlb[0].dns_name}"
    zone_id                = "${aws_lb.ingress-nlb[0].zone_id}"
    evaluate_target_health = true
  }
}
