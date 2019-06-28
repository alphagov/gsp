resource "aws_iam_role" "concourse" {
  name        = "${var.cluster_name}-concourse"
  description = "Role the concourse process assumes"

  assume_role_policy = "${data.aws_iam_policy_document.trust_kiam_server.json}"
}
