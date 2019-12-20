# ADR012: Docker image repositories

## Status

Superseded by [ADR028](ADR028-docker-tools.md)

## Context

Artefacts that we build must be published somewhere that a Kubernetes cluster
can pull from. This can be solved by using an image repository. Multiple options
exist including:

* [Docker Hub]
* Vendor Managed Registry ie. [AWS ECR][ECR] or [Google Cloud Container
  Registry][GCR]
* Self-hosted Docker Registry ie. [Docker Registry][Docker Registry Self Host]
  or [Harbor]

All the above implement the Docker registry API and are broadly equivalent.

Docker Hub will require managing authentication and authorisation separate to
the existing IAM setup.

Managed Registry from a vendor, would save us from dealing with authentication
or authorisation provided we are running the Kubernetes nodes with that vendor.

Running a self-hosted Docker registry will involved maintenance and operational
work from the team. But gives us a lot of flexibility.

As we are currently targeting AWS for our setup, it makes sense to delegate
extra work that way. It means, we can restrict access to the repositories with
the use of IAM policies.

## Decision

We will use Managed Registry from a vendor we are currently occupying, AWS ECR
at the time.

## Consequences

* Images must not contain sensitive data even if they are published to "private"
  registries.
* ECR requires creation of repository ahead of time, our build process will need
  to include that action.

[Docker Hub]: https://hub.docker.com/
[Docker Registry Self Host]: https://docs.docker.com/registry/deploying/
[ECR]: https://aws.amazon.com/ecr/
[GCR]: https://cloud.google.com/container-registry/
[Harbor]: https://github.com/goharbor/harbor
