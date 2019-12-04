# ADR044: Security improvements

## Status

Accepted

## Context

Currently there are no restrictions on what options a service team
puts into pod specifications.  This means that a pod in a tenant
namespace could escalate its privileges to access resources that it
should not have access to.

Examples include:

* host namespaces
* privileged containers
* host networking
* arbitrary host volume mounts

For a full list see [pod security policies][] in the kubernetes
documentation.

Currently, our defence against this is that we enforce deploying
infrastructure as code, and we require 2-eyes review on every code
change.  However, kubernetes offers pod security policies as an extra
layer of defence.

[pod security policies]: https://kubernetes.io/docs/concepts/policy/pod-security-policy/

## Decision

We will implement pod security policies according to the following:

### Cluster defaults

```
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: gsp-default
  annotations:
    seccomp.security.alpha.kubernetes.io/allowedProfileNames: '*'
spec:
  privileged: false
  allowPrivilegeEscalation: true
  allowedCapabilities:
  - '*'
  # Allow core volume types.
  volumes:
    - 'configMap'
    - 'emptyDir'
    - 'projected'
    - 'secret'
    - 'downwardAPI'
    # Assume that persistentVolumes set up by the cluster admin are safe to use.
    - 'persistentVolumeClaim'
  hostNetwork: false
  hostIPC: false
  hostPID: false
  runAsUser:
    rule: 'RunAsAny'
  seLinux:
    rule: 'RunAsAny'
  supplementalGroups:
    rule: 'RunAsAny'
  fsGroup:
    rule: 'RunAsAny'
  readOnlyRootFilesystem: false
```

### System namespaces

System namespaces, including `gsp-system`, `kube-system` and `istio-system`,
should use a less restrictive policy.

```
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: system
  annotations:
    seccomp.security.alpha.kubernetes.io/allowedProfileNames: '*'
spec:
  privileged: true
  allowPrivilegeEscalation: true
  allowedCapabilities:
  - '*'
  volumes:
  - '*'
  hostNetwork: true
  hostPorts:
  - min: 0
    max: 65535
  hostIPC: true
  hostPID: true
  runAsUser:
    rule: 'RunAsAny'
  seLinux:
    rule: 'RunAsAny'
  supplementalGroups:
    rule: 'RunAsAny'
  fsGroup:
    rule: 'RunAsAny'
```

## Consequences

System services will be able to get the access they need (privileged, host
networking etc.) but tenant service teams will be restricted in what they are
able to do. Tenant service teams will still be able to run container processes
as root (to be addressed in future if we so desire) but will not be able to
escalate to access parts of the host.
