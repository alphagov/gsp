# Writing deployment charts

## Overview

Your deployments are declared as a [Helm Chart](https://docs.helm.sh/developing_charts/).

A Chart is a collection of one or more [Kubernetes Objects](https://kubernetes.io/docs/concepts/#kubernetes-objects) (YAML definitions of [Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/), [Services](https://kubernetes.io/docs/concepts/services-networking/service/), etc).

Unlike raw Kubernetes Object YAML, a Chart allows some basic templating features to encourage reusability. For example you may have an [Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/) definition where you want to parameterize a domain like:

```
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: name-virtual-host-ingress
spec:
  rules:
  - host: myapp.{{ .Values.cluster.domain }}
    http:
      paths:
      - backend:
          serviceName: service2
          servicePort: 80
```

You can read more about the templating system in the [developing charts docs](https://docs.helm.sh/developing_charts/).

## Global values

The follow values are available for use in your Charts to help you write more reusable deployments.

| name | description | example |
|---|---|---|
| `{{ .Values.cluster.name }}` | the name of the environment | `"staging"` |
| `{{ .Values.cluster.domain }}` | the public dns zone of the cluster | `"run-sandbox.aws.ext.govsvc.uk"` |

