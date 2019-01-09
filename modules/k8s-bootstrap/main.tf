module "common" {
  source = "../common-ignition"
}

module "k8s-node" {
  source                = "../k8s-node-ignition"
  dns_service_ip        = "${var.dns_service_ip}"
  node_labels           = "node-role.kubernetes.io/master,node-role.kubernetes.io/bootstrapper"
  node_taints           = "node-role.kubernetes.io/master=:NoSchedule"
  cluster_domain_suffix = "${var.cluster_domain_suffix}"
  k8s_tag               = "${var.k8s_tag}"
  kubelet_kubeconfig    = "${var.kubelet_kubeconfig}"
  kube_ca_crt           = "${var.kube_ca_crt}"
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

resource "aws_instance" "bootstrap" {
  ami                         = "${data.aws_ami.coreos.id}"
  instance_type               = "${var.instance_type}"
  vpc_security_group_ids      = ["${var.security_group_ids}"]
  subnet_id                   = "${var.subnet_id}"
  associate_public_ip_address = "false"
  user_data                   = "${data.ignition_config.bootstrap-actual.rendered}"

  tags = "${map(
        "Name", "${var.cluster_name}-bootstrap",
        "kubernetes.io/cluster/${var.cluster_name}", "1",
        "KubernetesCluster", "${var.cluster_name}",
        "kubernetes.io/role/master", "1"
    )}"

  root_block_device = {
    volume_type = "gp2"
    volume_size = "40"
  }

  iam_instance_profile = "${var.iam_instance_profile_name}"
}

resource "aws_lb_target_group_attachment" "bootstrap" {
  target_group_arn = "${var.lb_target_group_arn}"
  target_id        = "${aws_instance.bootstrap.id}"
  port             = 6443
}
