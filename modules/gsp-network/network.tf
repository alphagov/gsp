resource "aws_vpc" "network" {
  cidr_block           = "${var.host_cidr}"
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = "${map(
    "Name", "${var.cluster_name}",
    "kubernetes.io/cluster/${var.cluster_name}", "shared",
  )}"
}

resource "aws_internet_gateway" "gateway" {
  vpc_id = "${aws_vpc.network.id}"

  tags = "${map("Name", "${var.cluster_name}")}"
}

locals {
  public_cidr_block  = "${cidrsubnet(aws_vpc.network.cidr_block, 1, 0)}"
  private_cidr_block = "${cidrsubnet(aws_vpc.network.cidr_block, 1, 1)}"
}

module "subnet-0" {
  source              = "../gsp-subnet"
  vpc_id              = "${aws_vpc.network.id}"
  cluster_name        = "${var.cluster_name}"
  private_cidr_block  = "${cidrsubnet(local.private_cidr_block, ceil(log(6, 2)), 0)}"
  public_cidr_block   = "${cidrsubnet(local.public_cidr_block, ceil(log(6, 2)), 0)}"
  availability_zone   = "eu-west-2a"
  internet_gateway_id = "${aws_internet_gateway.gateway.id}"
}

module "subnet-1" {
  source              = "../gsp-subnet"
  vpc_id              = "${aws_vpc.network.id}"
  cluster_name        = "${var.cluster_name}"
  private_cidr_block  = "${cidrsubnet(local.private_cidr_block, ceil(log(6, 2)), 1)}"
  public_cidr_block   = "${cidrsubnet(local.public_cidr_block, ceil(log(6, 2)), 1)}"
  availability_zone   = "eu-west-2b"
  internet_gateway_id = "${aws_internet_gateway.gateway.id}"
}

module "subnet-2" {
  source              = "../gsp-subnet"
  vpc_id              = "${aws_vpc.network.id}"
  cluster_name        = "${var.cluster_name}"
  private_cidr_block  = "${cidrsubnet(local.private_cidr_block, ceil(log(6, 2)), 2)}"
  public_cidr_block   = "${cidrsubnet(local.public_cidr_block, ceil(log(6, 2)), 2)}"
  availability_zone   = "eu-west-2c"
  internet_gateway_id = "${aws_internet_gateway.gateway.id}"
}
