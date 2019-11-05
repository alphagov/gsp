resource "aws_subnet" "public" {
  vpc_id                  = var.vpc_id
  availability_zone       = var.availability_zone
  cidr_block              = var.public_cidr_block
  map_public_ip_on_launch = false

  tags = {
    "Name"                                      = "${var.cluster_name}-public-${var.availability_zone}"
    "kubernetes.io/cluster/${var.cluster_name}" = "shared"
    "kubernetes.io/role/elb"                    = "1"
  }
}

resource "aws_route_table" "public" {
  vpc_id = var.vpc_id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = var.internet_gateway_id
  }

  tags = {
    "Name" = "${var.cluster_name}-public-${var.availability_zone}"
  }
}

resource "aws_route_table_association" "public" {
  route_table_id = aws_route_table.public.id
  subnet_id      = aws_subnet.public.id
}

resource "aws_eip" "egress" {
  vpc = "true"

  tags = {
    "Name" = "${var.cluster_name}-egress-${var.availability_zone}"
  }

  lifecycle {
    prevent_destroy = true
  }
}

resource "aws_nat_gateway" "egress" {
  allocation_id = aws_eip.egress.id
  subnet_id     = aws_subnet.public.id

  tags = {
    "Name" = "${var.cluster_name}-egress-${var.availability_zone}"
  }
}

resource "aws_eip" "ingress" {
  vpc = "true"

  tags = {
    "Name" = "${var.cluster_name}-ingress-${var.availability_zone}"
  }

  lifecycle {
    prevent_destroy = true
  }
}

