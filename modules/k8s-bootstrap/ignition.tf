data "template_file" "bootstrap-user-data-object-key" {
  template = "/user_data/$${cluster_name}-bootstrap.userdata"

  vars {
    cluster_name = "${var.cluster_name}"
  }
}

data "ignition_config" "bootstrap" {
  files = ["${module.k8s-node.ignition_file_ids}"]

  systemd = ["${concat(
        module.k8s-node.ignition_systemd_unit_ids,
        module.common.ignition-systemd-unit-ids,
    )}"]

  append = {
    source       = "${var.bootstrap_base_userdata_source}"
    verification = "${var.bootstrap_base_userdata_verification}"
  }
}

resource "aws_s3_bucket_object" "bootstrap-user-data" {
  bucket  = "${var.user_data_bucket_name}"
  key     = "${data.template_file.bootstrap-user-data-object-key.rendered}"
  content = "${data.ignition_config.bootstrap.rendered}"
}

data "ignition_config" "bootstrap-actual" {
  replace = {
    source       = "s3://${var.user_data_bucket_name}${aws_s3_bucket_object.bootstrap-user-data.id}"
    verification = "sha512-${sha512(data.ignition_config.bootstrap.rendered)}"
  }
}
