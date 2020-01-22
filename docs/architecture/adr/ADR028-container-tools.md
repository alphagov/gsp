# ADR028: Container Tools

## Status

Superseded by [ADR046](ADR046-replace-harbor-ecr.md)

## Context

We are currently using Docker as our container runtime.

There are needs for public docker images:

* so that master builds can be easily pulled and tested in the local development environments
* so that images can be easily shared between different teams

There are needs for digitally verifying the provenance of images:

* so that production systems can trust that an image has not been tampered with
* so that production systems can authenticate the origin of a build

There are needs for vulnerability scanning:

* so that production systems can warn or prevent exploitable software running in production

The docker ecosystem provides tooling that can help us meet these needs:

* [Docker Content Trust] (Notary) can be used to sign images and prove provenance
* [Docker Registries][Docker Registry] can expose images publicly
* Scanning tools like [Clair] can periodically or at pull/push time perform CVE scanning.

Unfortunately AWS [ECR] does not currently support public images or [Docker Content Trust], and there is no managed solution to image scanning from AWS as yet.

### Option 1: Wrap AWS ECR

We could write/manage a proxy to allow exposing [ECR] publicly and integrate the missing features.

* Potentially fragile implementation tied to the underlying AWS services
* Another thing to have to maintain
* Might offer ability to remove functionality as AWS support more features in future

### Option 2: Use an external SaaS offering

Use a SaaS service like [Quay] which offers most of these features.

* Additional configuration for cluster to pull from external source
* Reduces ability to automate provisioning (requires additional credential management)
* Procurement

### Option 3: Self-hosted Docker Tools in cluster

We could deploy Docker Distribution, Notary & Clair into the cluster backed by a managed storage backend like S3

* Well integrated with the platform
* Would work for local GSP instance


## Decision

We will run a self hosted set of Docker tools

## Consequences

* More to manage

[Clair]: https://github.com/coreos/clair
[Docker Hub]: https://hub.docker.com/
[Docker Content Trust]: https://docs.docker.com/engine/security/trust/content_trust/
[Docker Registry]: https://docs.docker.com/registry/
[Docker Registry API]: https://docs.docker.com/registry/spec/api/
[Docker Registry Self Host]: https://docs.docker.com/registry/deploying/
[ECR]: https://aws.amazon.com/ecr/
[GCR]: https://cloud.google.com/container-registry/
[Harbor]: https://github.com/goharbor/harbor
[Quay]: https://quay.io/
[S3]: https:/aws.amazon.com/S3
