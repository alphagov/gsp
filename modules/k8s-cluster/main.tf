data "aws_caller_identity" "current" {}

data "aws_region" "current" {}

data "aws_subnet" "private_subnets" {
  count = "${length(var.private_subnet_ids)}"
  id    = "${element(var.private_subnet_ids, count.index)}"
}

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

data "aws_iam_policy_document" "worker-nodes-assume-role-policy" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals = {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "worker-nodes-role" {
  name               = "${var.cluster_name}-worker-nodes-role"
  assume_role_policy = "${data.aws_iam_policy_document.worker-nodes-assume-role-policy.json}"
}

resource "aws_iam_instance_profile" "worker-nodes-profile" {
  name = "${var.cluster_name}-worker-nodes-profile"
  role = "${aws_iam_role.worker-nodes-role.name}"
}

resource "aws_iam_role_policy_attachment" "worker-nodes-eks-worker-policy-attachment" {
  role       = "${aws_iam_role.worker-nodes-role.name}"
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
}

resource "aws_iam_role_policy_attachment" "worker-nodes-eks-cni-policy-attachment" {
  role       = "${aws_iam_role.worker-nodes-role.name}"
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
}

resource "aws_iam_role_policy_attachment" "worker-nodes-ecr-ro-policy-attachment" {
  role       = "${aws_iam_role.worker-nodes-role.name}"
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}

resource "aws_lb_target_group" "worker-nodes-http-target-group" {
  name     = "worker-nodes-http-target-group"
  port     = 31380
  protocol = "HTTP"
  vpc_id   = "${var.vpc_id}"
}

resource "aws_lb_target_group" "worker-nodes-tcp-target-group" {
  name     = "worker-nodes-tcp-target-group"
  port     = 31390
  protocol = "TCP"
  vpc_id   = "${var.vpc_id}"
}

# As per https://docs.aws.amazon.com/eks/latest/userguide/launch-workers.html
resource "aws_cloudformation_stack" "worker-nodes-per-az" {
  count         = "${length(var.private_subnet_ids)}"
  name          = "${var.cluster_name}-worker-nodes-${element(data.aws_subnet.private_subnets.*.availability_zone, count.index)}"
  template_body = "${file("${path.module}/data/nodegroup-v2.yaml")}"
  capabilities  = ["CAPABILITY_IAM"]

  parameters = {
    ClusterName                      = "${var.cluster_name}"
    ClusterControlPlaneSecurityGroup = "${aws_security_group.controller.id}"
    NodeGroupName                    = "worker-${element(data.aws_subnet.private_subnets.*.availability_zone, count.index)}"

    NodeAutoScalingGroupMinSize         = "${var.minimum_workers_per_az_count}"
    NodeAutoScalingGroupDesiredCapacity = "${var.minimum_workers_per_az_count}"
    NodeAutoScalingGroupMaxSize         = "${var.maximum_workers_per_az_count}"

    NodeInstanceType    = "${var.worker_instance_type}"
    NodeInstanceProfile = "${aws_iam_instance_profile.worker-nodes-profile.arn}"
    NodeVolumeSize      = "40"
    BootstrapArguments  = "--kubelet-extra-args \"--node-labels=node-role.kubernetes.io/worker --event-qps=0\""
    VpcId               = "${var.vpc_id}"
    Subnets             = "${element(data.aws_subnet.private_subnets.*.id, count.index)}"
    NodeSecurityGroups  = "${aws_security_group.node.id},${aws_security_group.worker.id}"
    NodeTargetGroups    = "${aws_lb_target_group.worker-nodes-http-target-group.arn},${aws_lb_target_group.worker-nodes-tcp-target-group.arn}"
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
