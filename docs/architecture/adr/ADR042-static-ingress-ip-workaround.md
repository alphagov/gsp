# ADR042: Static ingress IP workaround

## Status

Accepted

## Context

Some production systems have existing requirements for static ingress IPs and
have external users/companies relying on being able to safelist ingress IPs.

The current ingress architecture is that each namespace within the cluster has
it's own dynamically provisioned NLB (Network Load Balancer) that directs traffic to
the namespace's own ingressgateway. [See
ADR037](./ADR037-per-namespace-gateways.md) for more details on this
architecture.

Static IP addresses create complications for dynamic, ephemeral infrastructure
systems such as this and as such we would like to continue to push back on
these kinds of static configurations that add little to security at a great
cost to complexity. However in the meantime we need a way to allocate static
ingress IPs.

We have already provisioned an advertised 5x potential ingress IPs to the
relevant external parties. But these EIPs have not yet been allocated.

It is not possible to attach EIPs to an existing NLB without recreating it.

Some research in this space:

* [Merged into Kubernetes v1.16.0](https://github.com/kubernetes/kubernetes/pull/69263/commits/7767535426b29fc14461083528b0d15493a3262e)
  is a change that allows allocating EIPs to dynamically provisioned LoadBalancer
  Services. This would allow us to provision static AWS Elastic IPs and attach
  them to the ingress load balancers as we require. Unfortunately at time of
  writing EKS is at `v1.14` and so we do not expect to reach `v1.16` for around 6 months.
* [AWS Global Accelerator](https://aws.amazon.com/global-accelerator/) is a
  cross-region load-balancer service that provides static IPs and allows
  forwarding traffic to regional-load-balancers or directly to instance.

Options:

1. We could manually provision a Global Accelerator to point at the
   cluster-provisioned ingress NLB.
1. We could manually provision a new NLB to point at all worker nodes on
   designated NodePort
1. ~~We could terraform provision a Global Accelerator (or new NLB) to point at the
   cluster-provisioned ingress NLB.~~ not supported in terraform at this time.
1. We could create a controller to manage the provisioning of Global
   Accelerator to point at the cluster-provisioned ingress NLB based on custom
   and/or Service resource configuration
1. We could create a controller to manage the provisioning of a new NLB to
   point at all worker nodes on a designated NodePort based on custom and/or Service
   resource configuration

## Decision

We will manually provision a Global Accelerator pointing to the ingress NLB
using the AWS console to point at the cluster-provisioned ingress NLB. **This
is intended to be a temporary solution until we are running a version of the
EKS control plane that supports attaching EIPs to LoadBalancer Services**

We choose the Global Accelerator (rather than a second load-balancer) as the
traffic shaping features offer a good route to deprecating this workaround in
the future and it is less complicated that attempting to duplicate the required
TargetGroups for the dynamically provisioned worker nodes.

We choose not to write a controller to manage the provisioning due to the
expectation that this will be a temporary solution and we do not expect anyone
else to benefit from the engineering effort.

## Consequences

* Since the ingress NLBs are dynamically managed by EKS, there is potential for
  the NLB to get "disconnected" from the manually provisioned endpoints.
* The workaround will only affect the designated namespace(s) and will
  potentially lead to differences between "staging/testing" and "production"
  deployment configurations.

## Appendix

Notes from a brief spike into this configuration can be found
[in the Global Accelerator Spike notes](../notes/global-accelerator-spike.md)
