data "template_file" "user-data-object-key" {
  template = "/user_data/$${cluster_name}-bootstrap-base.userdata"

  vars {
    cluster_name = "${var.cluster_name}"
  }
}

data "ignition_config" "bootstrap" {
  files = ["${concat(
        module.bootkube-assets.ignition-file-ids,
        module.etcd-cluster.bootkube-ignition-file-ids,
    )}"]

  systemd = ["${module.bootkube-assets.bootkube_systemd_unit_id}"]
}

resource "aws_s3_bucket_object" "bootstrap-user-data" {
  bucket  = "${var.user_data_bucket_name}"
  key     = "${data.template_file.user-data-object-key.rendered}"
  content = "${data.ignition_config.bootstrap.rendered}"
}
