# ADR014: Sealed Secrets

## Status

Accepted

## Context

We want to provide a simple way to pass sensitive values into environments via git.

Currently the only way to do this is by directly interacting with the cluster to inject secrets.

## Decision

We will deploy a [SealedSecrets](https://github.com/bitnami-labs/sealed-secrets) controller that allows sealing (encrypting) Kubernetes Secrets with a public key unique to each environment, making them safe to store as part of their deployment.

## Consequences

* A per-cluster key pair will require that secrets will have to be "sealed" once per cluster/environment.
* A per-cluster key pair will prevent teamA from ever being able to decrypt teamB's secrets by compromising their own cluster.
* A per-cluster key pair will reduce blast area of a leaked key.
* A per-cluster key pair may lead to more complex deployment charts that require selecting secrets based on environment.
