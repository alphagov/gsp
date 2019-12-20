# ADR007: Identity Provider

## Status

Superseded by [ADR023](ADR023-cluster-authentication.md)

## Context

We need to provide a way to authenticate users who will interact with our Kubernetes clusters.

We do not have a organisation-wide identity provider. Virtually everyone will have a Google account. Many people will have a GitHub account.

People working on GitHub repositories are likely the same people who are deploying to a cluster. Access to repositories likely indicates which users should have access to a cluster. We can reuse this user:team mapping in order to control access to clusters.

## Decision

We will use GitHub as our identify provider.

## Consequences

- It may be non-trivial for non-technical people to authenticate if they have to create a GitHub account first.
- Misconfiguration (accidental or malicious) in GitHub of users and organisations will allow/disallow cluster access.
- Granularity of permissions is limited by GitHub's permission model.
