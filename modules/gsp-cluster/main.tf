module "k8s-cluster" {
  source                       = "../k8s-cluster"
  vpc_id                       = "${var.vpc_id}"
  private_subnet_ids           = ["${var.private_subnet_ids}"]
  public_subnet_ids            = ["${var.public_subnet_ids}"]
  cluster_name                 = "${var.cluster_name}"
  worker_count                 = "${var.worker_count}"
  worker_instance_type         = "${var.worker_instance_type}"
  minimum_workers_per_az_count = "${var.minimum_workers_per_az_count}"
  maximum_workers_per_az_count = "${var.maximum_workers_per_az_count}"
  ci_worker_count              = "${var.ci_worker_count}"
  ci_worker_instance_type      = "${var.ci_worker_instance_type}"
  eks_version                  = "${var.eks_version}"

  apiserver_allowed_cidrs = ["${concat(
      formatlist("%s/32", var.egress_ips),
      var.gds_external_cidrs,
  )}"]
}
