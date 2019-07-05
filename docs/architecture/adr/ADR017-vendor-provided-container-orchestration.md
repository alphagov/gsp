# ADR017: Vendor-provided Container Orchestration

## Status

Pending

## Context

Following the [rollout of AWS EKS in London](https://aws.amazon.com/about-aws/whats-new/2019/02/amazon-eks-available-in-mumbai--london--and-paris-aws-regions/) it is now an attractive alternative to the hand-rolled kubernetes installation that was created as a result of [ADR003](ADR003-container-orchestration.md). This will bring numerous benefits:
* reducing the amount of infrastructure we manage in-house, by offloading it to AWS
* better alignment with [Technology & Operations Strategic Principle #3 - "Use fully managed cloud services by default"](https://reliability-engineering.cloudapps.digital/documentation/strategy-and-principles/re-principles.html#3-use-fully-managed-cloud-services-by-default)

As of 1.12, EKS supports what we need (e.g. Istio, kiam etc.).

## Decision

We will host the GDS Supported Platform on AWS EKS.

## Consequences

* We can no longer make changes to the control plane configuration
* Resources like DNS records and load balancers now have to be managed from within the cluster
* Reduce the cost of hosting the control plane
