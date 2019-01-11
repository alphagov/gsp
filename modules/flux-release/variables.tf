variable "namespace" {
    description = "namespace to deploy into"
    type = "string"
}

variable "chart_git" {
    description = "git repository containing helm chart to watch/deploy"
    type = "string"
}

variable "chart_ref" {
    description = "git ref/branch to watch"
    type = "string"
    default = "master"
}

variable "chart_path" {
    description = "path within the git repository to a helm chart to deploy"
    type = "string"
    default = ""
}

variable "addons_dir" {
    description = "local target path to place kubernetes resource yaml"
    type = "string"
    default = "addons"
}

variable "values" {
    description = "embedded yaml to pass to the helm resource for flux helm operator. Whitespace is important"
    type = "string"
    default = ""
}

variable "valueFileSecrets" {
    description = "List of names of Secrets containing additional helm values"
    type = "list"
    default = []
}

variable "cluster_name" {
    description = "name of this cluster/environment (accessible as .Values.cluster.name in charts)"
    type = "string"
}

variable "cluster_domain" {
    description = "domain mapped to this cluster/environment (accessible as .Values.cluster.name in charts)"
    type = "string"
}
