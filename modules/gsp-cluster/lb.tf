resource "aws_lb" "lb" {
  name                             = "${var.cluster_name}-lb"
  internal                         = "false"
  load_balancer_type               = "network"
  subnets                          = ["${aws_subnet.cluster-public.*.id}"]
  enable_cross_zone_load_balancing = "true"
}

resource "aws_lb_target_group" "controllers" {
  name        = "${var.cluster_name}-controllers"
  vpc_id      = "${aws_vpc.network.id}"
  target_type = "instance"

  protocol = "TCP"
  port     = 6443

  health_check {
    protocol = "TCP"
    port     = 6443

    healthy_threshold   = 3
    unhealthy_threshold = 3

    interval = 10
  }
}

resource "aws_lb_listener" "apiserver-https" {
  load_balancer_arn = "${aws_lb.lb.arn}"
  protocol          = "TCP"
  port              = "6443"

  default_action {
    type             = "forward"
    target_group_arn = "${aws_lb_target_group.controllers.arn}"
  }
}

resource "aws_lb_target_group" "workers-http" {
  name        = "${var.cluster_name}-workers-http"
  vpc_id      = "${aws_vpc.network.id}"
  target_type = "instance"

  protocol = "TCP"
  port     = 80

  health_check {
    protocol = "HTTP"
    port     = 10254
    path     = "/healthz"

    healthy_threshold   = 3
    unhealthy_threshold = 3

    interval = 10
  }
}

resource "aws_lb_listener" "ingress-http" {
  load_balancer_arn = "${aws_lb.lb.arn}"
  protocol          = "TCP"
  port              = 80

  default_action {
    type             = "forward"
    target_group_arn = "${aws_lb_target_group.workers-http.arn}"
  }
}

resource "aws_lb_target_group" "workers-https" {
  name        = "${var.cluster_name}-workers-https"
  vpc_id      = "${aws_vpc.network.id}"
  target_type = "instance"

  protocol = "TCP"
  port     = 443

  health_check {
    protocol = "HTTP"
    port     = 10254
    path     = "/healthz"

    healthy_threshold   = 3
    unhealthy_threshold = 3

    interval = 10
  }
}

resource "aws_lb_listener" "ingress-https" {
  load_balancer_arn = "${aws_lb.lb.arn}"
  protocol          = "TCP"
  port              = 443

  default_action {
    type             = "forward"
    target_group_arn = "${aws_lb_target_group.workers-https.arn}"
  }
}

resource "aws_route53_record" "apiserver" {
  zone_id = "${var.dns_zone_id}"

  name = "${format("%s.%s.", var.cluster_name, var.dns_zone)}"
  type = "A"

  alias {
    name                   = "${aws_lb.lb.dns_name}"
    zone_id                = "${aws_lb.lb.zone_id}"
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "ingress" {
  zone_id = "${var.dns_zone_id}"
  name    = "*.${var.cluster_name}.${var.dns_zone}"
  type    = "CNAME"
  ttl     = "300"
  records = ["${aws_lb.lb.dns_name}"]
}
