# ADR006: Cluster Authentication Method

## Status

Superseded by [ADR023](ADR023-cluster-authentication.md)

## Context

[Kubernetes offers a bunch of strategies for authentication](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#authentication-strategies). The most appropriate strategies are x509 certificates, webhook tokens and OpenID Connect (OIDC) Tokens.

Using x509 certificates would allow us to authenticate to a cluster without having to pass additional `--oidc-*` flags to the Kubernetes masters. This is something that is often hard or impossible to change on managed clusters. However, it isn't clear that managed clusters that we'd be able authenticate using certificates anyway. For example, EKS seems to only allow IAM for authentication.

Renewal of short-lived x509 certificates would require regenerating, revalidating and resigning a certificate as there is no equivalent refresh token as with [OIDC tokens](https://auth0.com/docs/tokens/refresh-token/current). This is not part of `kubectl`.

Webhook tokens require additional configuration in the apiserver (`--authentication-token-webhook-*` flags and a configuration file). They allow a supplied token to be validated by an external service and kubernetes apiserver expects a certain payload in the response containing the user properties. An application to generate a token for a given user (following authentication with an identity provider such as GitHub) will need to be written. If expiry of tokens is a requirement this will also need to be managed. [Guard](https://github.com/appscode/guard) was looked at in evaluating flows for webhook tokens.

OIDC tokens require passing `--oidc-*` arguments to the Kubernetes apiservers. However, the refreshing of short-lived tokens is [natively supported by `kubectl`](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#using-kubectl).

## Decision

We will use OIDC.

Webhooks force a permanent inclusion of additional configuration, potentially additional application code to be maintained and don't offer the flexibility and security of OIDC tokens. X509 certificates come with all the complexities of managing certificates but offer no additional security or flexibility over OIDC.

## Consequences

- We will have to run and manage our own clusters until it becomes possible to specify `--oidc-*` arguments on managed masters.
- We will have to run an OIDC client.
- We may have to run an OIDC provider (such as [Dex](https://github.com/dexidp/dex)) if we choose to authenticate with a provider that does not support OIDC natively (such as GitHub).
