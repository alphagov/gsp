# ADR037: Per namespace istio gateways

## Status

Accepted

## Context

Different tenant applications have differing ingress and egress requirements. Examples include:

* mutual TLS
* egress controls (e.g. limiting what can egress from a namespace / cluster)
* certificate operations (renewal etc.)

There is also the vision of the GSP that provides a common platform for all tenants based off a shared codebase. This platform allows suitable customisations to be suitable for the tenants running on it and, where possible, those customisations are controlled via kube-yaml.

At present, a single ingress gateway and egress gateway are defined by the GSP. In the case of the ingress gateway, all ingress traffic is routed through this which creates a range of difficulties:

* secrets relating to TLS certs (mutual or otherwise) have to be stored in `istio-system`
* having some endpoints use mutual TLS and others using non-mutual TLS is difficult/awkward/impossible on a single gateway
  * a gateway is able to load certificates from a single secret (by default) so different applications in different namespaces will inevitably conflict
  * different tenants in the same cluster will need to make changes to common "system" parts of the code (i.e. `istio-system`)
* controlling DNS entries and TLS configuration per tenant is difficult/awkward/impossible
  * DNS is currently controlled via terraform in the GSP codebase (a single "*." record for the domain for the cluster provided as a convenience).
  * TLS is currently terminated outside the cluter in an Application Load Balancer (ALB) in AWS, which doesn't suit all tenants

Other, indirect features:

* istio routing rules for single gateways in `istio-system` are:
  * difficult to reason about and become quite brittle as a result
  * don't align very well with istio's documentation

Options

1. "Snowflake" load balancers and DNS records in tenant-specific terraform, bypass the `istio-system` ingress gateway as needed

  Pros:

  * Quick
  * Relatively easy & familar

  Cons:

  * Custom terraform for a tenant for something likely to be common (e.g. managing TLS termination and DNS ingress)
  * Doesn't make istio routing easier to work with
  * Doesn't improve certificate management (provisioning, renewing etc.)

1. Allow tenants to configure multiple secrets on the `istio-system` ingress gateways

  Pros:
  * Addresses the single-secret certificate problem

  Cons:

  * Tenants need to know about gateway configuration
  * Tenants still manipulating shared namespaces & resources
  * Doesn't make istio routing easier to work with
  * Doesn't improve certificate management (provisioning, renewing etc.)

1. Cert-manager and external-dns with single `istio-system` gateways

  Pros:

  * Certificate management
  * DNS management (potentially)

  Cons:

  * Tenants still manipulating shared namespaces & resources
  * Doesn't make istio routing easier to work with

1. Per-namespace gateways with cert-manager and external-dns

  Pros:

  * Certificate management
  * DNS management
  * Ingress control
  * Improves istio routing management (as it aligns more with the docs)
  * Alows tenants direct control (hence it's more inline with the vision)

  Cons:

  * A complex change

## Decision

We will create an ingress gateway and an egress gateway in each tenant ("managed") namespace (option 4 above).

## Consequences

* tenants will have control over ingress rules (TLS, mutual TLS, secrets, routing, certificates, DNS)
* an increase in route53 records
* an increase in AWS load balancer provisions
* may create some downtime during the initial migration
