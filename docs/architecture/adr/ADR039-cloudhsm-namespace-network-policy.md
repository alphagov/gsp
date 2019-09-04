# ADR039: Restricting CloudHSM network access to particular namespaces

## Status

Accepted

## Context

[ADR036](ADR036-hsm-isolation-in-detail.md) described the network and
credential isolation we use to ensure that unauthorised users cannot
access the CloudHSM.

Recently in 3ea9de2ff, we introduced a GlobalNetworkPolicy object,
which is a Calico feature that allows a cluster-wide network policy to
be imposed.  This allows us to control network access to and from
particular namespaces in a way which cannot be overridden by tenants.

In particular, currently access to the HSM is only allowed from pods
annotated with a `talksToHsm=true` label.

When working in their own namespace, a developer has full control over
what labels they put on their pods, so they can still choose to put
the `talksToHsm=true` label on their pods.  But they do not have
control over what labels the namespace itself has; to change this
would require a change to the `gsp` or appropriate `cluster-config`
repository, which would make such a change visible to many more
people.

Therefore, if we extend the GlobalNetworkPolicy to require a
`talksToHsm=true` label on *both* the pod *and* the namespace, we will
prevent tenants from unilaterally opening up network access to the HSM
from their namespaces.

## Decision

We will augment the GlobalNetworkPolicy (previously described in ADR036) by:

 - setting a `GlobalNetworkPolicy` that denies access to the
   CloudHSM's IP address unless the pod carries a label
   (`talksToHsm=true`) and the namespace also carries a label
   (`talksToHsm=true`) and allows all other egress traffic

## Consequences

Control of which namespaces get the `talksToHsm=true` label will be
via the appropriate `-cluster-config` repo.  If a developer wants to
allow a new namespace to talk to the HSM, they will need to issue a PR
against that repo.

If we are confident in the GlobalNetworkPolicy's control of HSM
access, we could consider reducing the technical controls required on
non-HSM namespaces.  For example, we could consider allowing
developers to run plain `kubectl apply` in unprivileged namespaces,
for fast-feedback learning.
