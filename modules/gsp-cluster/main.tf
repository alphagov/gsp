module "k8s-cluster" {
  source = "../k8s-cluster"
  vpc_id = var.vpc_id

  private_subnet_ids = var.private_subnet_ids
  public_subnet_ids  = var.public_subnet_ids
  cluster_name       = var.cluster_name

  minimum_workers_per_az_count           = var.minimum_workers_per_az_count
  desired_workers_per_az_map             = var.desired_workers_per_az_map
  maximum_workers_per_az_count           = var.maximum_workers_per_az_count
  worker_on_demand_base_capacity         = var.worker_on_demand_base_capacity
  worker_on_demand_percentage_above_base = var.worker_on_demand_percentage_above_base

  eks_version        = var.eks_version
  worker_eks_version = var.worker_eks_version
  apiserver_allowed_cidrs = concat(
    formatlist("%s/32", var.egress_ips),
    var.gds_external_cidrs,
  )
  worker_generation_timestamp = var.worker_generation_timestamp
}

