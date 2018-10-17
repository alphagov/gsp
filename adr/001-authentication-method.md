# ADR 001: Cluster Authentication Method

## Status

PROPOSED

## Context

[Kubernetes offers a bunch of strategies for authentication](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#authentication-strategies). The two most appropriate strategies are x509 certificates and OpenID Connect (OIDC) Tokens.

Using x509 certificates would allow us to authenticate to a cluster without having to pass additional `--oidc-*` flags to the Kubernetes masters. This is something that is often hard or impossible to change on managed clusters. However, it isn't clear that managed clusters that we'd be able authenticate using certificates anyway. For example, EKS seems to only allow IAM for authentication.

Renewal of short-lived x509 certificates would require regenerating, revalidating and resigning a certificate as there is no equivalent refresh token as with [OIDC tokens](https://auth0.com/docs/tokens/refresh-token/current). This is not part of `kubectl`.

OIDC tokens require modifying `--oidc-*` arguments to Kubernetes masters. However, the refreshing of short-lived tokens is [natively supported by `kubectl`](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#using-kubectl).

## Decision

We have chosen to use OIDC.

## Consequences

- We will have to run and manage our own clusters until it becomes possible to specify `--oidc-*` arguments on managed masters.
- We will have to run an OIDC client.
- We may have to run an OIDC provider (such as [Dex](https://github.com/dexidp/dex)) if we choose to authenticate with a provider that does not support OIDC natively (such as GitHub).
