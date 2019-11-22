resource "aws_iam_role" "concourse" {
  name        = "${var.cluster_name}-concourse"
  description = "Role the concourse process assumes"

  assume_role_policy = data.aws_iam_policy_document.trust_kiam_server.json
}

resource "random_password" "concourse_db_master_password" {
  length  = 32
  special = false
}

resource "aws_rds_cluster" "concourse" {
  cluster_identifier = "${var.cluster_name}-concourse"
  database_name      = "concourse"
  master_password    = random_password.concourse_db_master_password.result
  master_username    = "concourse"
  availability_zones = var.availability_zones
  engine             = "aurora-postgresql"
  engine_version     = "10.7"
  enabled_cloudwatch_logs_exports = [
    "audit",
    "error",
    "general",
    "slowquery",
    "postgresql",
  ]
  db_subnet_group_name   = aws_db_subnet_group.private.name
  vpc_security_group_ids = [aws_security_group.rds-from-worker.id]
  tags = {
    Name = var.cluster_name
  }
}

resource "aws_rds_cluster_instance" "concourse" {
  count                = 2
  identifier           = "${var.cluster_name}-concourse-${count.index}"
  cluster_identifier   = "${var.cluster_name}-concourse"
  instance_class       = "db.t3.medium"
  db_subnet_group_name = aws_db_subnet_group.private.name
  tags = {
    Name = var.cluster_name
  }
}
