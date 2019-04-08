resource "local_file" "auth-delegator" {
    content  = "${file("${path.module}/data/metrics-server/auth-delegator.yaml")}"
    filename = "addons/${var.cluster_name}/metrics-server/auth-delegator.yaml"
}

resource "local_file" "auth-reader" {
    content  = "${file("${path.module}/data/metrics-server/auth-reader.yaml")}"
    filename = "addons/${var.cluster_name}/metrics-server/auth-reader.yaml"
}


resource "local_file" "metrics-apiservice" {
    content  = "${file("${path.module}/data/metrics-server/metrics-apiservice.yaml")}"
    filename = "addons/${var.cluster_name}/metrics-server/metrics-apiservice.yaml"
}


resource "local_file" "metrics-server-deployment" {
    content  = "${file("${path.module}/data/metrics-server/metrics-server-deployment.yaml")}"
    filename = "addons/${var.cluster_name}/metrics-server/metrics-server-deployment.yaml"
}


resource "local_file" "metrics-server-service" {
    content  = "${file("${path.module}/data/metrics-server/metrics-server-service.yaml")}"
    filename = "addons/${var.cluster_name}/metrics-server/metrics-server-service.yaml"
}


resource "local_file" "resource-reader" {
    content  = "${file("${path.module}/data/metrics-server/resource-reader.yaml")}"
    filename = "addons/${var.cluster_name}/metrics-server/resource-reader.yaml"
}
