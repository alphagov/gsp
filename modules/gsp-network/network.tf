resource "aws_vpc" "network" {
  cidr_block           = "10.${var.netnum}.0.0/16"
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

module "subnet-0" {
  source              = "../gsp-subnet"
  vpc_id              = "${aws_vpc.network.id}"
  cluster_name        = "${var.cluster_name}"
  availability_zone   = "eu-west-2a"
  internet_gateway_id = "${aws_internet_gateway.gateway.id}"
  private_cidr_block  = "10.${var.netnum}.0.0/19"
  public_cidr_block   = "10.${var.netnum}.32.0/20"

  # protected_cidr_block = "10.${var.netnum}.48.0/21"
  # spare_cidr_block     = "10.${var.netnum}.56.0/21"
}

module "subnet-1" {
  source              = "../gsp-subnet"
  vpc_id              = "${aws_vpc.network.id}"
  cluster_name        = "${var.cluster_name}"
  availability_zone   = "eu-west-2b"
  internet_gateway_id = "${aws_internet_gateway.gateway.id}"
  private_cidr_block  = "10.${var.netnum}.64.0/19"
  public_cidr_block   = "10.${var.netnum}.96.0/20"

  # protected_cidr_block = "10.${var.netnum}.112.0/21"
  # spare_cidr_block     = "10.${var.netnum}.120.0/21"
}

module "subnet-2" {
  source              = "../gsp-subnet"
  vpc_id              = "${aws_vpc.network.id}"
  cluster_name        = "${var.cluster_name}"
  availability_zone   = "eu-west-2c"
  internet_gateway_id = "${aws_internet_gateway.gateway.id}"
  private_cidr_block  = "10.${var.netnum}.128.0/19"
  public_cidr_block   = "10.${var.netnum}.160.0/20"

  # protected_cidr_block = "10.${var.netnum}.176.0/21"
  # spare_cidr_block     = "10.${var.netnum}.184.0/21"
}

# following range is left over for future needs
# spare_subnet_block = "10.${var.netnum}.192.0/18"

