# ADR004: Tenant Isolation

## Status

Superseded by [ADR024 - Soft Multitenancy](ADR024-soft-multitenancy.md)

## Context

The two main isolation models for Kubernetes are:

* Namespaces within a single cluster
* Running multiple clusters

All Service Teams currently have separate AWS accounts.

Some Service Teams have separate AWS accounts for separate environment (ie. Staging, Production etc)

Many Service Teams have micro-service architectures

Some Service Teams have unique network isolation requirements that may be hard to implement in a shared environment.

To ensure "smooth transition" during a migration it would be preferable to have clusters deployed to Service Team's VPCs.

To ensure separation of billing it would be preferable to deploy clusters to Service Team's AWS accounts.

To ensure strong network/compute isolation between Service Teams it would be preferable to deploy separate clusters for separate environments.

## Decision

We will manage multiple Kubernetes clusters, one per environment per tenant.

This should provide:

* Strong network/compute isolation between service teams
* Clear separation of billing
* Ability to deploy into target VPCs to make piecemeal microservice migration practical

## Consequences

* Management multiple clusters across multiple accounts may be more complex than managing a single big one
* Lack of direct control over the target account/environment may introduce complexities
