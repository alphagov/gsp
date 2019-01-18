data "aws_availability_zones" "all" {}

resource "aws_vpc" "network" {
  cidr_block           = "${var.host_cidr}"
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = "${map("Name", "${var.cluster_name}")}"
}

resource "aws_internet_gateway" "gateway" {
  vpc_id = "${aws_vpc.network.id}"

  tags = "${map("Name", "${var.cluster_name}")}"
}

resource "aws_route_table" "cluster-public" {
  vpc_id = "${aws_vpc.network.id}"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.gateway.id}"
  }

  tags = "${map("Name", "${var.cluster_name}")}"
}

resource "aws_route_table" "cluster-private" {
  vpc_id = "${aws_vpc.network.id}"

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = "${element(aws_nat_gateway.cluster.*.id, count.index)}"
  }

  tags = "${map("Name", "${var.cluster_name}")}"
}

resource "aws_subnet" "cluster-private" {
  count = "${length(data.aws_availability_zones.all.names)}"

  vpc_id            = "${aws_vpc.network.id}"
  availability_zone = "${data.aws_availability_zones.all.names[count.index]}"

  cidr_block              = "${cidrsubnet(var.host_cidr, 4, count.index)}"
  map_public_ip_on_launch = false

  tags = "${map("Name", "${var.cluster_name}-cluster-${count.index}")}"
}

resource "aws_subnet" "cluster-public" {
  count = "${length(data.aws_availability_zones.all.names)}"

  vpc_id            = "${aws_vpc.network.id}"
  availability_zone = "${data.aws_availability_zones.all.names[count.index]}"

  cidr_block              = "${cidrsubnet(var.host_cidr, 4, count.index + length(data.aws_availability_zones.all.names))}"
  map_public_ip_on_launch = false

  tags = "${map("Name", "${var.cluster_name}-cluster-${count.index}")}"
}

resource "aws_route_table_association" "cluster-private" {
  count = "${length(data.aws_availability_zones.all.names)}"

  route_table_id = "${element(aws_route_table.cluster-private.*.id, count.index)}"
  subnet_id      = "${element(aws_subnet.cluster-private.*.id, count.index)}"
}

resource "aws_route_table_association" "cluster-public" {
  count = "${length(data.aws_availability_zones.all.names)}"

  route_table_id = "${element(aws_route_table.cluster-public.*.id, count.index)}"
  subnet_id      = "${element(aws_subnet.cluster-public.*.id, count.index)}"
}

resource "aws_eip" "public" {
  count = "${length(data.aws_availability_zones.all.names)}"

  vpc = "true"
}

resource "aws_nat_gateway" "cluster" {
  count = "${length(data.aws_availability_zones.all.names)}"

  allocation_id = "${element(aws_eip.public.*.id, count.index)}"
  subnet_id     = "${element(aws_subnet.cluster-public.*.id, count.index)}"
  depends_on    = ["aws_internet_gateway.gateway"]
}
