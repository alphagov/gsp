# ADR046: Replacing harbor with ECR

## Status

Accepted

## Context

[Harbor][] has proved itself to be rather unstable (causing several incidents from
various root causes). We have also been gradually rolling back our usage within
GSP of the various tools that Harbor provides:

* we are not using notary for signing (image provenance); this is now being
  handled using image digests instead of tags, enforced by gatekeeper
* we haven't yet implemented anything using Clair
* we no longer have a requirement for anonymous public access of docker repositories

It is also worth noting that some of these (e.g. "Clair") have been integrated
into AWS ECR since we last looked.

In order to implement the decisions in [ADR045 dev namespaces][ADR045] we must
ensure that each namespace has its own place to store docker images, and that
one namespace cannot overwrite another namespace's docker images. This would be
very fiddly to achieve with Harbor because it would require interactions with
the Harbor API which doesn't fit very well with "desired state configuration".
However with ECR we can more easily leverage IAM and the existing Service
Operator patterns to allow tenants to provision their own ECR repositories with
their own credentials.

See also [ADR012][].

## Decision

We will replace the Harbor registry (the last Harbor component in use in GSP at
time of writing) with AWS ECR with "push" credentials managed by the [Service
Operator][].

## Consequences

We will be able to deliver namespace-isolated repositories for pushing images,
private repositories only and stability improvements to the GSP as a whole by
removing Harbor completely. This will also mean that the existence of a
repository in ECR is managed by yaml.

[ADR012]: ADR012-docker-image-repositories.md
[ADR045]: ADR045-dev-namespaces.md
[Service Operator]:
https://github.com/alphagov/gsp/tree/master/components/service-operator
[Harbor]: https://goharbor.io/
