# Helm Chart Best Practices

## Overview

Guidance for authoring Helm charts suitable for production deployments.

### Use the release name in resource names

`.metadata.name` should include at least `{{ .Release.Name }}`

For example:

```
metadata:
    name: {{ .Release.Name }}-my-app
```

This allows you to install multiple versions of your application to a single
namespace without clashing/clobbering resources with the same name.

### High Availability

Replicas should be specified per application, rather than globally so that they can we scaled individually:

```
myApp:
  replicas: 2
```

### Resources & Limits

resources should be specified in pod templates (both requests and limits) when parameterised in values.yaml files like:

```
  myApp:
    resources:
      requests:
        cpu: 500m
        memory: 128Mi
      limits:
        cpu: 1000m
        memory: 512Mi
```

Which can then be consumed in templates with the `toYaml` helper:

```
  spec:
    template:
      spec:
        containers:
          - name: my-app
            resources:
  {{ toYaml .Values.resources | indent 12 }}
```

### Readiness and liveness probes

Use both readiness probes and liveness probes

### Set an update strategy on deployments

```
  spec:
    strategy:
      type: RollingUpdate
      rollingUpdate:
        maxUnavailable: 25%
        maxSurge: 25%
  ```

### One resource per file

Avoid `kind: List` resources and aim for a single resource per file.

## Use `enabled` convention

Sometimes it's useful to turn on/off features of a deployment, use a boolean flag in values.yaml defaults:

```
  stubs:
    enabled: true
```

...and guard entire files of configuration with an if block to include/exclude it:

```
  {{- if .Values.stubs.enabled }}
  ...
  {{- end }}
```

