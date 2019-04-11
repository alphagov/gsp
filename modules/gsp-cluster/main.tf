module "k8s-cluster" {
  source                  = "../k8s-cluster"
  vpc_id                  = "${var.vpc_id}"
  subnet_ids              = ["${concat(var.private_subnet_ids, var.public_subnet_ids)}"]
  cluster_name            = "${var.cluster_name}"
  worker_count            = "${var.worker_count}"
  worker_instance_type    = "${var.worker_instance_type}"
  ci_worker_count         = "${var.ci_worker_count}"
  ci_worker_instance_type = "${var.ci_worker_instance_type}"
  eks_version             = "${var.eks_version}"

  apiserver_allowed_cidrs = ["${concat(
      formatlist("%s/32", var.nat_gateway_public_ips),
      var.gds_external_cidrs,
  )}"]
}
