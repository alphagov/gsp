data "template_file" "namespace" {
  template = "${file("${path.module}/data/namespace.yaml")}"

  vars {
    permitted_roles_regex         = "${var.permitted_roles_regex}"
    namespace                     = "${var.namespace}"
    extra_namespace_labels        = "${var.extra_namespace_labels}"
  }
}

resource "local_file" "namespace" {
  count    = "${var.enabled == 0 ? 0 : 1}"
  filename = "${var.addons_dir}/${var.namespace}-namespace.yaml"
  content  = "${data.template_file.namespace.rendered}"
}

data "template_file" "helm-release" {
  template = "${file("${path.module}/data/helm-release.yaml")}"

  vars {
    namespace        = "${var.namespace}"
    release_name     = "${coalesce(var.release_name, var.namespace)}"
    chart_git        = "${var.chart_git}"
    chart_ref        = "${var.chart_ref}"
    chart_path       = "${var.chart_path}"
    cluster_name     = "${var.cluster_name}"
    cluster_domain   = "${var.cluster_domain}"
    values           = "${var.values}"
    valueFileSecrets = "[${join(",",formatlist("{\"name\":\"%s\"}", var.valueFileSecrets))}]"
    verification_keys = "[${join(",",formatlist("%#v", var.verification_keys))}]"
  }
}

resource "local_file" "helm-release-yaml" {
  count    = "${var.enabled == 0 ? 0 : 1}"
  filename = "${var.addons_dir}/${var.namespace}-helm-release.yaml"
  content  = "${data.template_file.helm-release.rendered}"
}

data "template_file" "values" {
  count    = "${length(var.valueFileSecrets)}"
  template = "${file("${path.module}/data/values-secret.yaml")}"

  vars {
    namespace = "${var.namespace}"
    name      = "${var.valueFileSecrets[count.index]}"
  }
}

resource "local_file" "values" {
  count    = "${length(var.valueFileSecrets)}"
  filename = "${var.addons_dir}/${var.namespace}-${var.valueFileSecrets[count.index]}-values.yaml"
  content  = "${data.template_file.values.*.rendered[count.index]}"
}

resource "tls_private_key" "github_deployment_key" {
  count     = "${var.enabled == 0 ? 0 : 1}"
  algorithm = "RSA"
  rsa_bits  = "4096"
}

resource "local_file" "ci-secrets" {
  count    = "${var.enabled == 0 ? 0 : 1}"
  filename = "${var.addons_dir}/${var.cluster_name}/${var.namespace}-deploy-keys.yaml"
  content  = "${data.template_file.ci-secrets.rendered}"
}

data "template_file" "ci-secrets" {
  count    = "${var.enabled == 0 ? 0 : 1}"
  template = "${file("${path.module}/data/ci-deploy-keys.yaml")}"

  vars = {
    ci_namespace = "ci-system-main" # The namespace for the secrets to be read by concourse. Should be `ci-system-[TEAM_NAME]`.
    namespace    = "${var.namespace}" # The namespace of the actual team's release deployment.
    private_key  = "${base64encode(element(concat(tls_private_key.github_deployment_key.*.private_key_pem, list("")), count.index))}"
    public_key   = "${element(concat(tls_private_key.github_deployment_key.*.public_key_openssh, list("")), count.index)}"
    secret_name  = "${var.namespace}"
  }
}

output "release-name" {
  value = "${coalesce(var.release_name, var.namespace)}"
}
