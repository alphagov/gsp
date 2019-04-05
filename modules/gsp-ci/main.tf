resource "aws_s3_bucket" "ci-system-harbor-registry-storage" {
  count = "${var.enabled == 0 ? 0 : 1}"

  bucket = "registry-${var.cluster_name}-${replace(var.dns_zone, ".", "-")}"
  acl    = "private"

  force_destroy = true # NEED TO VALIDATE!!!

  tags = {
    Name = "Harbor registry and chartmuseum storage"
  }
}

resource "local_file" "concourse-web-configmap" {
  count    = "${var.enabled == 0 ? 0 : 1}"
  filename = "addons/${var.cluster_name}/${module.ci-system.release-name}-web-configmap.yaml"
  content  = "${data.template_file.concourse-web-configmap.rendered}"
}

data "template_file" "concourse-web-configmap" {
  count    = "${var.enabled == 0 ? 0 : 1}"
  template = "${file("${path.module}/data/web-configmap.yaml")}"

  vars = {
    namespace    = "ci-system"
    release_name = "${module.ci-system.release-name}"
    teams        = "${jsonencode(var.github_teams)}"
  }
}

module "ci-system" {
  source = "..//flux-release"

  enabled               = "${var.enabled == 0 ? 0 : 1}"
  namespace             = "ci-system"
  chart_git             = "https://github.com/alphagov/gsp-ci-system.git"
  chart_ref             = "8dde966d3795fe04d42ed969d126822c5134dfb2"
  cluster_name          = "${var.cluster_name}"
  cluster_domain        = "${var.cluster_name}.${var.dns_zone}"
  addons_dir            = "addons/${var.cluster_name}"
  permitted_roles_regex = "^${var.enabled ? element(concat(aws_iam_role.harbor.*.name, list("")), 0) : ""}$"

  values = <<EOF
    concourse:
      secrets:
        localUsers: "pipeline-operator:${random_string.concourse_password.result}"
        githubClientId: "${var.github_client_id}"
        githubClientSecret: "${var.github_client_secret}"
        githubCaCert: "${var.github_ca_cert}"
      web:
        additionalVolumes:
        - name: "${module.ci-system.release-name}-web-configuration"
          configMap:
            name: "${module.ci-system.release-name}-web-configuration"
        additionalVolumeMounts:
        - name: "${module.ci-system.release-name}-web-configuration"
          mountPath: /web-configuration
        ingress:
          enabled: true
          annotations:
            kubernetes.io/tls-acme: "true"
          hosts:
          - "ci.${var.cluster_name}.${var.dns_zone}"
          tls:
          - secretName: concourse-web-tls
            hosts:
            - "ci.${var.cluster_name}.${var.dns_zone}"
      concourse:
        web:
          externalUrl: "https://ci.${var.cluster_name}.${var.dns_zone}"
          auth:
            github:
              enabled: true
            mainTeam:
              localUser: "pipeline-operator"
              config: /web-configuration/config.yaml
          kubernetes:
            namespacePrefix: "${module.ci-system.release-name}-"
            createTeamNamespaces: false
            teams: ["${join(",", concat(list("main"), var.concourse_teams))}"]
    pipelineOperator:
      concourseUsername: "pipeline-operator"
      concoursePassword: "${random_string.concourse_password.result}"
    harbor:
      harborAdminPassword: "${random_string.harbor_password.result}"
      secretKey: "${random_string.harbor_secret_key.result}"
      externalURL: "https://registry.${var.cluster_name}.${var.dns_zone}"
      persistence:
        imageChartStorage:
          type: s3
          s3:
            bucket: ${var.enabled ? element(concat(aws_s3_bucket.ci-system-harbor-registry-storage.*.id, list("")), 0) : ""}
            region: ${var.enabled ? element(concat(aws_s3_bucket.ci-system-harbor-registry-storage.*.region, list("")), 0) : ""}
            regionendpoint: s3.${var.enabled ? element(concat(aws_s3_bucket.ci-system-harbor-registry-storage.*.region, list("")), 0) : ""}.amazonaws.com
      expose:
        tls:
          secretName: harbor-registry-certificates
          notarySecretName: harbor-notary-certificates
        ingress:
          annotations:
            kubernetes.io/tls-acme: "true"
          hosts:
            core: "registry.${var.cluster_name}.${var.dns_zone}"
            notary: "notary.${var.cluster_name}.${var.dns_zone}"
      registry:
        podAnnotations:
          iam.amazonaws.com/role: "${var.enabled ? element(concat(aws_iam_role.harbor.*.name, list("")), 0) : ""}"
      chartmuseum:
        podAnnotations:
          iam.amazonaws.com/role: "${var.enabled ? element(concat(aws_iam_role.harbor.*.name, list("")), 0) : ""}"
EOF
}
