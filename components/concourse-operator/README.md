# Concourse Operator

## Overview

Declaratively configure concourse pipelines as Kubernetes resources.

* Defines a Pipeline CRD to describe a concourse pipeline
* Deploys a concourse-pipeline-controller to create/update/delete pipelines on change

## Example Pipeline

Define a pipeline like:

```yaml
apiVersion: concourse.govsvc.uk/v1beta1
kind: Pipeline
metadata:
  labels:
    controller-tools.k8s.io: "1.0"
  name: canary
spec:
  config:
    resources:
    - name: gsp-canary
      type: git
      source:
        uri: https://github.com/alphagov/gsp-canary.git
        branch: master
    - name: canary-image
      type: docker-image
      source:
        repository: govsvc/gsp-canary
    - name: updater-image
      type: docker-image
      source:
        repository: govsvc/gsp-canary-chart-updater
    jobs:
    - name: build
      plan:
      - get: gsp-canary
      - put: canary-image
        params:
          build: gsp-canary
          build_args:
            BUILD_TIMESTAMP: "1544635666"
          dockerfile: gsp-canary/Dockerfile.canary
      - put: updater-image
        params:
          dockerfile: gsp-canary/Dockerfile.updater
```

## Example Team

Define a team like:

```yaml

apiVersion: concourse.govsvc.uk/v1beta1
kind: Team
metadata:
  labels:
    controller-tools.k8s.io: "1.0"
  name: team-sample
spec:
  roles:
  - name: owner
    github:
      users: ["admin"]
  - name: member
    github:
      teams: ["org:team"]
  - name: viewer
    github:
      orgs: ["org"]
    local:
      users: ["visitor"]
```

## Example operator deployment

The operator is deployed as part of https://github.com/alphagov/gsp.

## Developing

Driven by `make`. To build the docker image:

```
make docker-build
```

To tag and publish the built docker image:

```
make docker-push
```

It was originally built using [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) with:

```
kubebuilder init --domain k8s.io
kubebuilder create api --group concourse --version v1beta1 --kind Pipeline
```
