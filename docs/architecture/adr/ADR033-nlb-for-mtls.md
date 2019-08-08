# ADR033: NLB for mTLS

## Status

Accepted

## Context

Verify's [doc-checking service](https://github.com/alphagov/doc-checking) is
secured in part using mTLS. Currently, our clusters are fronted by ALBs which
cannot provide mTLS.

The doc-checking service currently runs an nginx that provides the mTLS
functionality. In order for GSP to be able to allow something within the
cluster to perform mTLS we must run a load balancer that forwards unaltered TCP
packets in addition to, or instead of an ALB.

## Decision

We will optionally create and run an NLB in addition to the current ALB for
clusters that have a requirement to terminate their own TLS. This NLB will be
available at `nlb.$CLUSTER_DOMAIN`.

## Consequences

In some of our clusters we will be running two public load balancers. This may
be confusing or unexpected.

The certificates in use for mTLS will be managed outside of ACM. This means any
certificates will have to manually rotated unless we decide to start running
`cert-manager` again.
