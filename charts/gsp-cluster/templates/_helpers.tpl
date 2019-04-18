{{- define "pipelineOperator.serviceAccountName" -}}
{{- printf "%s-%s" .Release.Name .Values.pipelineOperator.serviceAccountName -}}
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
