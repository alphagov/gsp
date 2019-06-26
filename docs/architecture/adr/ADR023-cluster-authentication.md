# ADR023: Cluster Authentication

## Status

Accepted

## Context

We need to provide a secure way for users to authenticate to interact with cluster resources.

There are four different roles identified based on need:


| Cluster Role | Need |
|---|---|
| deployer | ability to make changes to cluster and full access to AWS resources (for CI) |
| admin | ability to make changes to cluster resources, and restricted access to AWS resources |
| sre | read only access to all cluster resources |
| dev | read only access to resources potentially scoped to a namespace |

## Decision

We will authenticate all users to IAM roles via the [aws-iam-authenticator](https://github.com/kubernetes-sigs/aws-iam-authenticator) and map those IAM roles to [ClusterRoles](https://kubernetes.io/docs/reference/access-authn-authz/rbac/) within the GSP cluster.

We will store the mapping of IAM user ARN to Cluster Role in Github so that it can be verified. [gds-trusted-developers](https://github.com/alphagov/gds-trusted-developers)

## Consequences

* Requires all users to have an assumable IAM user
* Requires all users to install the aws-iam-authenticator binary to use kubectl
