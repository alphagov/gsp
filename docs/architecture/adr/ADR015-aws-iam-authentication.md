# ADR015: AWS IAM Authentication (for admins)

## Status

Superceded by [ADR032](ADR032-aws-iam-authentication.md)

## Context

IAM Roles that can be assumed by authorized infrastructure engineers currently do not give access to the clusters via kubectl. We do not want to have to manage two sets of admins.

## Decision

We will enable any admin-like roles within the cluster only to those who can authenicate via the [aws-iam-authenticator](https://github.com/kubernetes-sigs/aws-iam-authenticator) assuming an appropriate role within the AWS account.

This should provide:

* Auditing to CloudTrial of authentication attempts and more
* Single place to manage roles
* A way to enable MFA/policy for cluster access
* A better user-experience for accessing clusters (with help from aws-vault)
* Simpler distribution of cluster configs (kubeconfig not containing anything sensitive)

## Consequences

* Requires installation of additional binary
