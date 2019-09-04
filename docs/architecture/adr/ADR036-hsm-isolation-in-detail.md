# ADR036: CloudHSM isolation

## Status

Accepted, amended by [ADR039](ADR039-cloudhsm-namespace-network-policy.md)

## Context

Some of our GSP clusters provide a CloudHSM instance.  It is important to
ensure that only authorised applications can use the CloudHSM.  There are
different choices available for technical controls to prevent unauthorised
access.

This ADR is a post hoc document; it documents existing practice and writes down
what we have already agreed and been doing for a while, but didn't capture as
an ADR.

### Network isolation

Network isolation ensures that only authorised applications are permitted to
open connections to the CloudHSM.  In a traditional AWS world, this could be
achieved using security groups - each app would have its own security group,
and the HSM security group would only allow connections from authorised apps.

In Kubernetes, all pods share the same subnet and same set of IPs.  This means
that the HSM security group cannot distinguish between traffic coming from
different apps.

Kubernetes provides [Network Policies][] to place controls on which pods are
allowed to connect to which other pods.  Network Policies are similar to
security groups, in that you can think of them as layer 4 stateful firewalls.

We have also deployed [Istio][], an advanced networking tool which provides
many features, but in the context of this ADR the main things it provides are
further network controls at a higher level (broadly, at layer 7).

The provider used for implementing the kubernetes network policies, [Calico],
extends the kubernetes API with more flexible resources, allowing for finer
grained control over routing rules. Calico's resources overlap with the native
kubernetes network policies and Istio's more powerful, higher-level traffic
management features.

[Network Policies]: https://kubernetes.io/docs/concepts/services-networking/network-policies/
[Istio]: https://istio.io
[Calico]: https://www.projectcalico.org/

### Authentication and credential management

In order to use the CloudHSM, you need to connect to it, but you also need to
authenticate to it.  CloudHSM supports multiple users.  Each user has a
different username and password.

The authorised eIDAS apps have a username and password provided to them in a
`Secret`.  Apps in other namespaces do not have access to this `Secret` and
therefore do not have the required username and password to authenticate to the
CloudHSM.

## Decision

We will prevent unauthorised network access to the CloudHSM by:

 - setting a `GlobalNetworkPolicy` that denies access to the CloudHSM's IP
   address unless the pod carries a label (`talksToHsm=true`) and allows all
   other egress traffic
 - setting each tenant namespace to use Istio's sidecar injection (the "service
   mesh" or just "mesh")
 - enabling Istio's mesh to only route traffic to destinations in its service
   registry (requiring an Istio `ServiceEntry` to be explicitly created in order
   to allow egress traffic from the mesh)
 
We will prevent unauthorised apps from having CloudHSM credentials by:

 - deploying the credentials in a `Secret` in the namespace for authorised apps
   only

## Consequences

A user that can create pods in the Verify cluster in a namespace other than the
eIDAS namespace will, provided they set the correct label, be able to:

 - establish a TCP connections to the CloudHSM used by eIDAS
 - run a subset of commands on the CloudHSM that [require do not require
   authentication](https://docs.aws.amazon.com/cloudhsm/latest/userguide/cloudhsm_mgmt_util-reference.html).
   Currently this appears to be the following commands:
  - `getHSMInfo`
  - `getKeyInfo`
  - `info`
  - `listUsers`
  - `quit`

A user that can somehow bypass the Istio-injected sidecar will be able to
establish TCP/UDP connections to any endpoint, but not the CloudHSM (unless it's
done from a pod with the correct label).
