module "common" {
  source = "../common-ignition"
}

data "aws_ami" "coreos" {
  most_recent = true
  owners      = ["595879546273"]

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "name"
    values = ["CoreOS-stable-*"]
  }
}

resource "aws_instance" "etcds" {
  count = "${var.node_count}"

  tags = "${map(
    "Name", "${var.cluster_name}-etcd${count.index}",
  )}"

  instance_type = "${var.instance_type}"

  ami       = "${data.aws_ami.coreos.image_id}"
  user_data = "${element(data.ignition_config.etcd-actual.*.rendered, count.index)}"

  root_block_device {
    volume_type = "gp2"
    volume_size = "40"
  }

  associate_public_ip_address = false
  subnet_id                   = "${element(var.subnet_ids, count.index % length(var.subnet_ids))}"
  vpc_security_group_ids      = ["${aws_security_group.etcd.id}"]

  lifecycle {
    ignore_changes = [
      "ami",
      "user_data",
    ]
  }

  iam_instance_profile = "${aws_iam_instance_profile.etcd_profile.name}"
}

data "template_file" "etcd_servers" {
  count    = "${var.node_count}"
  template = "$${cluster_name}-etcd$${index}.$${dns_zone}"

  vars {
    cluster_name = "${var.cluster_name}"
    index        = "${count.index}"
    dns_zone     = "${var.dns_zone}"
  }
}

resource "aws_route53_record" "etcds" {
  count = "${var.node_count}"

  zone_id = "${var.dns_zone_id}"

  name = "${format("%s.", element(data.template_file.etcd_servers.*.rendered, count.index))}"
  type = "A"
  ttl  = 300

  # private IPv4 address for etcd
  records = ["${element(aws_instance.etcds.*.private_ip, count.index)}"]
}
