# ADR047: Replace Kiam with IAM Roles for Service Accounts

## Status

Accepted

## Context

[Kiam][] was added into the GSP cluster in order to restrict IAM role access at a `Pod` and `Namespace`-level. It provides two things we care about. Kiam:

1. Allows an individual Pod to assume an IAM role without granting all Pods that ability.
1. Blocks access to the Instance Metadata Service.

Kiam achieves this by running "agents" on every node which proxy access to the [Instance Metadata Service][] (aka IMDS, aka `169.254.169.254`) through the Kiam "server". The Kiam server only proxies the request through to the "real" IMDS iff the `Pod` and `Namespace` have the correct annotations.

When we introduced Kiam there was no AWS-managed equivalent. Soon after we deployed Kiam to our clusters AWS released [IAM Roles for Service Accounts][].

We have experienced problems operating Kiam. Kiam breakages would often cause a Pod to be unable to assume an IAM Role that it required. These problems would often fail in fairly opaque way where it was not immediately clear that Kiam was not working as expected. Issues which caused downtime at some point include, but not limited to:

- Mismatch between the Kiam server and agent's mTLS certificates leading to communication failures.
- Kiam not being annotated with `priorityClassName: gsp-critical` leading to agents/servers not being scheduled/started.

Since we are running EKS we are able to use IAM Roles for Service Accounts to allow Pods to be given restricted access to IAM Roles. This could be used to provide the first feature that we require Kiam for.

In order to use IAM Roles for Service Accounts a [supported version of the AWS SDK][] must be used. We think the only thing we run that does not support IAM Roles for Service Accounts is [Grafana][]. [Whilst upstream is looking at making Grafana work with IAM Roles for Service Accounts][] it doesn't look like it will be available soon. In order for allow Grafana to continue to assume a role that has CloudWatch metrics access we will grant read-only access to CloudWatch metrics at an EC2 instance level.

The other feature that we rely on Kiam for is to block direct access to the IMDS. We will use a `GlobalNetworkPolicy` to make sure only correctly annotated Pods in `gsp-system` are able to access the IMDS.

## Decision

- Replace Kiam usage with an equivalent IAM role for service account.
- Grant all nodes read-only access to CloudWatch metrics so that Grafana can continue to access CloudWatch as a data source.
- Block all access to `169.254.169.254` from non-`gsp-system` namespaces.
- Only allow access to `169.254.169.254` from specifically annotated pods in the `gsp-system` namespace.
- Remove Kiam.

## Consequences

- Every node will have access to a role that has read-only access to a cluster's CloudWatch metrics. However, this shouldn't be accessible to any Pod outside of `gsp-system`.

[Kiam]: https://github.com/uswitch/kiam
[IAM Roles for Service Accounts]: https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html
[Grafana]: https://grafana.com/
[Instace Metadata Service]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html
[supported version of the AWS SDK]: https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts-minimum-sdk.html
[Whilst upstream is looking at making Grafana work with IAM Roles for Service Accounts]: https://github.com/grafana/grafana/issues/20473#issuecomment-559181587
