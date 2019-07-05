# ADR011: Build Artefacts

## Status

Pending

## Context

As part of our pipelines we will be building artefacts that will be used to test
and deploy our applications. We will be deploying applications to Kubernetes. We
will need to build a container image of some kind.

There are some competing container image formats, namely:

* [OCI]
* [ACI]

The OCI image format is [based on the Docker v2][oci-standard] image format.

The Kubernetes project appears to [prefer Docker/OCI][k8s-preferance] images
over ACI.

[rkt is moving to OCI][rkt-oci] and away from ACI. OCI will become the preferred
image format.

Docker has wide industry adoption and appears to have wide understanding within
GDS.

Docker is the default container runtime for Kubernetes.

## Decision

We will build and store OCI images built using Docker.

## Consequences

* It may be tricky to deploy apps outside of a container orchestrator of some
  kind.
* This means that the images will need to be pre-built and stored in accessible
  way.
* We will be unable to use container runtimes that do not support OCI images.

[OCI]: https://github.com/opencontainers/image-spec
[ACI]: https://github.com/appc/spec/blob/259c2eebc32df77c016974d5e8eed390d5a81500/spec/aci.md#app-container-image
[oci-standard]: https://blog.docker.com/2017/07/oci-release-of-v1-0-runtime-and-image-format-specifications/
[k8s-preferance]: https://kubernetes.io/blog/2015/05/docker-and-kubernetes-and-appc/
[rkt-oci]: https://github.com/rkt/rkt/blob/03285a7db960311faf887452538b2b8ae4304488/ROADMAP.md#oci-native-support
