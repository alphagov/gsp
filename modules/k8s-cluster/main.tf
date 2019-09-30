data "aws_caller_identity" "current" {}

data "aws_region" "current" {}

resource "aws_eks_cluster" "eks-cluster" {
  name     = "${var.cluster_name}"
  role_arn = "${aws_iam_role.eks-cluster.arn}"
  version  = "${var.eks_version}"

  vpc_config {
    security_group_ids = ["${aws_security_group.controller.id}"]
    subnet_ids         = ["${concat(var.private_subnet_ids, var.public_subnet_ids)}"]
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
    "aws_cloudwatch_log_group.eks",
  ]

  lifecycle {
    prevent_destroy = true
  }
}

resource "aws_iam_openid_connect_provider" "eks" {
  client_id_list  = ["sts.amazonaws.com"]
  thumbprint_list = []
  url             = "${aws_eks_cluster.eks-cluster.identity.0.oidc.0.issuer}"
}

resource "aws_cloudwatch_log_group" "eks" {
  name              = "/aws/eks/${var.cluster_name}/cluster"
  retention_in_days = 30

  lifecycle {
    prevent_destroy = true
  }
}

# As per https://docs.aws.amazon.com/eks/latest/userguide/launch-workers.html
resource "aws_cloudformation_stack" "worker-nodes" {
  name          = "${var.cluster_name}-worker-nodes"
  template_body = "${file("${path.module}/data/nodegroup.yaml")}"
  capabilities  = ["CAPABILITY_IAM"]

  parameters = {
    ClusterName                         = "${var.cluster_name}"
    ClusterControlPlaneSecurityGroup    = "${aws_security_group.controller.id}"
    NodeGroupName                       = "worker"
    NodeAutoScalingGroupMinSize         = "${var.worker_count}"
    NodeAutoScalingGroupDesiredCapacity = "${var.worker_count}"
    NodeAutoScalingGroupMaxSize         = "${var.worker_count + 2}"
    NodeInstanceType                    = "${var.worker_instance_type}"
    NodeVolumeSize                      = "40"
    BootstrapArguments                  = "--kubelet-extra-args \"--node-labels=node-role.kubernetes.io/worker --event-qps=0\""
    VpcId                               = "${var.vpc_id}"
    Subnets                             = "${join(",", var.private_subnet_ids)}"
  }

  depends_on = ["aws_eks_cluster.eks-cluster"]
}

resource "aws_cloudformation_stack" "kiam-server-nodes" {
  name          = "${var.cluster_name}-kiam-server-nodes"
  template_body = "${file("${path.module}/data/nodegroup.yaml")}"
  capabilities  = ["CAPABILITY_IAM"]

  parameters = {
    ClusterName                         = "${var.cluster_name}"
    ClusterControlPlaneSecurityGroup    = "${aws_security_group.controller.id}"
    NodeGroupName                       = "kiam"
    NodeAutoScalingGroupMinSize         = "2"
    NodeAutoScalingGroupDesiredCapacity = "2"
    NodeAutoScalingGroupMaxSize         = "3"
    NodeInstanceType                    = "t3.medium"
    NodeVolumeSize                      = "40"
    BootstrapArguments                  = "--kubelet-extra-args \"--node-labels=node-role.kubernetes.io/cluster-management --register-with-taints=node-role.kubernetes.io/cluster-management=:NoSchedule --event-qps=0\""
    VpcId                               = "${var.vpc_id}"
    Subnets                             = "${join(",", var.private_subnet_ids)}"
  }

  depends_on = ["aws_eks_cluster.eks-cluster"]
}

resource "aws_cloudformation_stack" "ci-nodes" {
  name          = "${var.cluster_name}-ci-nodes"
  template_body = "${file("${path.module}/data/nodegroup.yaml")}"
  capabilities  = ["CAPABILITY_IAM"]

  parameters = {
    ClusterName                         = "${var.cluster_name}"
    ClusterControlPlaneSecurityGroup    = "${aws_security_group.controller.id}"
    NodeGroupName                       = "ci"
    NodeAutoScalingGroupMinSize         = "${var.ci_worker_count}"
    NodeAutoScalingGroupDesiredCapacity = "${var.ci_worker_count}"
    NodeAutoScalingGroupMaxSize         = "${var.ci_worker_count + 1}"
    NodeInstanceType                    = "${var.ci_worker_instance_type}"
    NodeVolumeSize                      = "40"
    BootstrapArguments                  = "--kubelet-extra-args \"--node-labels=node-role.kubernetes.io/ci --register-with-taints=node-role.kubernetes.io/ci=:NoSchedule --event-qps=0\""
    VpcId                               = "${var.vpc_id}"
    Subnets                             = "${join(",", var.private_subnet_ids)}"
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
