{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "gsp-cluster.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "gsp-cluster.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "gsp-cluster.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}


{{- define "pipelineOperator.serviceAccountName" -}}
{{- printf "%s-%s" .Release.Name .Values.pipelineOperator.serviceAccountName -}}
{{- end -}}

{{- define "serviceOperator.serviceAccountName" -}}
{{- printf "%s-%s" .Release.Name .Values.serviceOperator.serviceAccountName -}}
{{- end -}}

{{/*
These concourse-related templates need to match the values
templated by the concourse chart - there doesn't appear to be
a way to grab the templated values out.
*/}}
{{- define "concourse.service.name" -}}
{{- printf "%s-%s" .Release.Name "web" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "concourse.namespace.prefix" -}}
{{- printf "%s-" .Release.Name -}}
{{- end -}}

{{- define "dockerconfig.creds" -}}
{{- printf "admin:%s" .Values.harbor.harborAdminPassword | b64enc -}}
{{- end -}}

{{- define "dockerconfig.json" -}}
{
    "auths": {
        "https://registry.{{ .Values.global.cluster.domain }}/": {
            "auth": "{{ include "dockerconfig.creds" . }}"
        }
    }
}
{{- end -}}
