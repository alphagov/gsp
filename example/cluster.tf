terraform {
  backend "s3" {
    bucket = "BUCKET_NAME"
    region = "REGION"
    key    = "STATEFILE_LOCATION"
  }
}

data "aws_caller_identity" "current" {}

module "gsp-cluster" {
    source = "git::https://github.com/alphagov/gsp-terraform-ignition//modules/gsp-cluster"
    cluster_name = "${var.cluster_name}"
    dns_zone = "${var.cluster_zone}"
    user_data_bucket_name = "${var.user_data_bucket_name}"
    user_data_bucket_region = "${var.user_data_bucket_region}"
    k8s_tag = "${var.k8s_tag}"
    admin_role_arns = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/admin"]
}
