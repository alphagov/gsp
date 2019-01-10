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

resource "aws_launch_configuration" "controller" {
  name_prefix       = "${var.cluster_name}-controller-"
  image_id          = "${data.aws_ami.coreos.image_id}"
  instance_type     = "${var.controller_instance_type}"
  enable_monitoring = false

  user_data = "${data.ignition_config.controller-actual.rendered}"

  root_block_device {
    volume_type = "gp2"
    volume_size = "40"
  }

  security_groups = ["${aws_security_group.controller.id}"]

  lifecycle {
    create_before_destroy = true
    ignore_changes = ["image_id"]
  }

  iam_instance_profile = "${aws_iam_instance_profile.controller_profile.name}"
}

resource "aws_autoscaling_group" "controllers" {
  name = "${var.cluster_name}-controller"

  desired_capacity          = "${var.controller_count}"
  min_size                  = "${var.controller_count}"
  max_size                  = "${var.controller_count}"
  default_cooldown          = 30
  health_check_grace_period = 30

  vpc_zone_identifier = ["${var.subnet_ids}"]

  launch_configuration = "${aws_launch_configuration.controller.name}"

  target_group_arns = ["${var.controller_target_group_arns}"]

  # Waiting for instance creation delays adding the ASG to state. If instances
  # can't be created (e.g. spot price too low), the ASG will be orphaned.
  # Orphaned ASGs escape cleanup, can't be updated, and keep bidding if spot is
  # used. Disable wait to avoid issues and align with other clouds.
  wait_for_capacity_timeout = "0"

  tags = [
    {
      key                 = "Name"
      value               = "${var.cluster_name}-controller"
      propagate_at_launch = true
    },
    {
      key                 = "KubernetesCluster"
      value               = "${var.cluster_name}"
      propagate_at_launch = true
    },
    {
      key                 = "kubernetes.io/cluster/${var.cluster_name}"
      value               = "1"
      propagate_at_launch = true
    },
    {
      key                 = "kubernetes.io/role/master"
      value               = "1"
      propagate_at_launch = true
    },
  ]
}

resource "aws_launch_configuration" "worker" {
  name_prefix       = "${var.cluster_name}-worker-"
  image_id          = "${data.aws_ami.coreos.image_id}"
  instance_type     = "${var.worker_instance_type}"
  enable_monitoring = false

  user_data = "${data.ignition_config.worker-actual.rendered}"

  root_block_device {
    volume_type = "gp2"
    volume_size = "40"
  }

  security_groups = ["${aws_security_group.worker.id}"]

  lifecycle {
    create_before_destroy = true
    ignore_changes        = ["image_id"]
  }

  iam_instance_profile = "${aws_iam_instance_profile.worker_profile.name}"
}

resource "aws_autoscaling_group" "workers" {
  name = "${var.cluster_name}-worker"

  desired_capacity          = "${var.worker_count}"
  min_size                  = "${var.worker_count}"
  max_size                  = "${var.worker_count}"
  default_cooldown          = 30
  health_check_grace_period = 30

  vpc_zone_identifier = ["${var.subnet_ids}"]

  launch_configuration = "${aws_launch_configuration.worker.name}"

  target_group_arns = ["${var.worker_target_group_arns}"]

  # Waiting for instance creation delays adding the ASG to state. If instances
  # can't be created (e.g. spot price too low), the ASG will be orphaned.
  # Orphaned ASGs escape cleanup, can't be updated, and keep bidding if spot is
  # used. Disable wait to avoid issues and align with other clouds.
  wait_for_capacity_timeout = "0"

  tags = [
    {
      key                 = "Name"
      value               = "${var.cluster_name}-worker"
      propagate_at_launch = true
    },
    {
      key                 = "KubernetesCluster"
      value               = "${var.cluster_name}"
      propagate_at_launch = true
    },
    {
      key                 = "kubernetes.io/cluster/${var.cluster_name}"
      value               = "1"
      propagate_at_launch = true
    },
  ]
}
