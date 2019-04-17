data "aws_caller_identity" "current" {}

data "aws_region" "current" {}

resource "aws_eks_cluster" "eks-cluster" {
  name     = "${var.cluster_name}"
  role_arn = "${aws_iam_role.eks-cluster.arn}"
  version  = "${var.eks_version}"

  vpc_config {
    security_group_ids = ["${aws_security_group.controller.id}"]
    subnet_ids         = ["${var.subnet_ids}"]
  }

  enabled_cluster_log_types = [
    "api",
    "audit",
    "authenticator",
    "controllerManager",
    "scheduler",
  ]

  depends_on = [
    "aws_iam_role_policy_attachment.eks-cluster-policy",
    "aws_iam_role_policy_attachment.eks-service-policy",
  ]
}

resource "aws_key_pair" "eks" {
  key_name   = "${var.cluster_name}"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDP1p65o186ezp0tpDid1qaIbNUg3QZmHqRtDiGOuzndebqK50uMKJ9KUGCPGXVEK+S1ztzMuRy/NB/U59UKxOyx/M6/oKclRH58+YRVg+ftp4hDNE9HW5o3mWvGmgOLtlaCfrpwXbOgtrT0pRH1R0qqeX0J3hY6m8JQlH+cHYdj7e/HjJQSpKaKmyBakQU8wHjvX4yjtxBRLdoLcOVQapHZPs2iFU8sqYT3FIGHSf3lyF/j+I9gNxe/B1KsTZpp+FrQ9mve7uSruK4fS0FvjTVZ/eemMm4niuAn4KGdFzmbHU3bBV1MpS04d4xrVQzhI/tXXVS8ZXF9Xekg2xtDmLWvl4hn43nfJa+8RSUVoeBv1XQRasafkgmF+ysTJhgsNZC8jziiEk7YqZh+uqkNRNIpveZU3bU+6aOzisM9VCA/HSoIEHLYUswXAXRbi75JMIhAT0l/RBTYpqnboOdb+MXM2jbcbJ2xbUJcFfDrpMHt1wE1Y+sp0P0zGEx9LubWQhplYYgMmW56NcYyHQTSS+V3EXaoEki17Qgg3MSMIdrWgH0c9EBnVj0L0dIvDzdYzUpbav0DAKu9ElYXcjnbDzZLmEyfZf5pnSC+NMCW2CF6E8C+FvcJ9akP+C3IU+tn5cc4u8eG5XuXh2vChXAs7B6slzvLPjpgbofuwVdErCrbw== daniel.blair@digital.cabinet-office.gov.uk"
}

# As per https://docs.aws.amazon.com/eks/latest/userguide/launch-workers.html
resource "aws_cloudformation_stack" "worker-nodes" {
  name         = "${var.cluster_name}-worker-nodes"
  template_url = "https://amazon-eks.s3-us-west-2.amazonaws.com/cloudformation/2019-02-11/amazon-eks-nodegroup.yaml"
  capabilities = ["CAPABILITY_IAM"]

  parameters = {
    ClusterName                         = "${var.cluster_name}"
    ClusterControlPlaneSecurityGroup    = "${aws_security_group.controller.id}"
    NodeGroupName                       = "${var.cluster_name}-worker-nodes"
    NodeAutoScalingGroupMinSize         = "${var.worker_count}"
    NodeAutoScalingGroupDesiredCapacity = "${var.worker_count}"
    NodeAutoScalingGroupMaxSize         = "${var.worker_count + 1}"
    NodeInstanceType                    = "${var.worker_instance_type}"
    NodeImageId                         = "ami-0c7388116d474ee10"
    NodeVolumeSize                      = "40"
    KeyName                             = "${aws_key_pair.eks.key_name}"
    BootstrapArguments                  = ""
    VpcId                               = "${var.vpc_id}"
    Subnets                             = "${join(",", var.subnet_ids)}"
  }

  depends_on = ["aws_eks_cluster.eks-cluster"]
}

resource "aws_cloudformation_stack" "kiam-server-nodes" {
  name         = "${var.cluster_name}-kiam-server-nodes"
  template_url = "https://amazon-eks.s3-us-west-2.amazonaws.com/cloudformation/2019-02-11/amazon-eks-nodegroup.yaml"
  capabilities = ["CAPABILITY_IAM"]

  parameters = {
    ClusterName                         = "${var.cluster_name}"
    ClusterControlPlaneSecurityGroup    = "${aws_security_group.controller.id}"
    NodeGroupName                       = "${var.cluster_name}-kiam-server-nodes"
    NodeAutoScalingGroupMinSize         = "2"
    NodeAutoScalingGroupDesiredCapacity = "2"
    NodeAutoScalingGroupMaxSize         = "3"
    NodeInstanceType                    = "t2.small"
    NodeImageId                         = "ami-0c7388116d474ee10"
    NodeVolumeSize                      = "40"
    KeyName                             = "${aws_key_pair.eks.key_name}"
    BootstrapArguments                  = "--kubelet-extra-args \"--node-labels=node-role.kubernetes.io/cluster-management --register-with-taints=node-role.kubernetes.io/cluster-management=:NoSchedule\""
    VpcId                               = "${var.vpc_id}"
    Subnets                             = "${join(",", var.subnet_ids)}"
  }

  depends_on = ["aws_eks_cluster.eks-cluster"]
}

resource "aws_cloudformation_stack" "ci-nodes" {
  name         = "${var.cluster_name}-ci-nodes"
  template_url = "https://amazon-eks.s3-us-west-2.amazonaws.com/cloudformation/2019-02-11/amazon-eks-nodegroup.yaml"
  capabilities = ["CAPABILITY_IAM"]

  parameters = {
    ClusterName                         = "${var.cluster_name}"
    ClusterControlPlaneSecurityGroup    = "${aws_security_group.controller.id}"
    NodeGroupName                       = "${var.cluster_name}-ci-nodes"
    NodeAutoScalingGroupMinSize         = "${var.ci_worker_count}"
    NodeAutoScalingGroupDesiredCapacity = "${var.ci_worker_count}"
    NodeAutoScalingGroupMaxSize         = "${var.ci_worker_count + 1}"
    NodeInstanceType                    = "${var.ci_worker_instance_type}"
    NodeImageId                         = "ami-0c7388116d474ee10"
    NodeVolumeSize                      = "40"
    KeyName                             = "${aws_key_pair.eks.key_name}"
    BootstrapArguments                  = "--kubelet-extra-args \"--node-labels=node-role.kubernetes.io/ci --register-with-taints=node-role.kubernetes.io/ci=:NoSchedule\""
    VpcId                               = "${var.vpc_id}"
    Subnets                             = "${join(",", var.subnet_ids)}"
  }

  depends_on = ["aws_eks_cluster.eks-cluster"]
}

data "template_file" "kubeconfig" {
  template = "${file("${path.module}/data/kubeconfig")}"

  vars {
    apiserver_endpoint = "${aws_eks_cluster.eks-cluster.endpoint}"
    ca_cert            = "${aws_eks_cluster.eks-cluster.certificate_authority.0.data}"
    name               = "${var.cluster_name}"
    cluster_id         = "${var.cluster_name}"
  }
}

data "template_file" "aws-auth" {
  template = "${file("${path.module}/data/aws-auth.yaml")}"

  vars {
    bootstrapper_role_mappings = "${join("\n", formatlist(var.bootstrapper_role_arn_mapping_template, list(aws_cloudformation_stack.worker-nodes.outputs["NodeInstanceRole"], aws_cloudformation_stack.kiam-server-nodes.outputs["NodeInstanceRole"], aws_cloudformation_stack.ci-nodes.outputs["NodeInstanceRole"])))}"
    iam_admin_role_mappings    = "${join("\n", formatlist(var.admin_role_arn_mapping_template, var.admin_role_arns))}"
    iam_sre_role_mappings      = "${join("\n", formatlist(var.sre_role_arn_mapping_template, var.sre_role_arns))}"
    iam_dev_role_mappings      = "${join("\n", formatlist(var.dev_role_arn_mapping_template, var.dev_role_arns))}"
  }
}
