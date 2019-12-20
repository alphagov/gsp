# ADR009: Multi-tenancy for CI and CD

## Status

Superseded by [ADR029](ADR029-continuous-delivery-tools.md)

## Context

Two models have been proposed concerning CI and CD tool sets:

1. Multi-tenant: all tenants, including Reliability Engineering, share that same CI and CD instance
2. Per-tenant: each tenant has their own CI and CD cluster

## Decision

- There will be a single CI and CD toolset used by all tenants of the new service

## Consequences

- With a single CI and CD toolset used by all tenants, there will be a need for those tools to have strong RBAC to prevent cross-tenant pollution
- The decision to use a single tool for all tenants could easily be adapted in the future to allow per-tenant CI and CD tool deployments should requirements change
- Authentication for the CI and CD toolsets should align with [ADR007](ADR007-identity-provider.md)
- It is believed that by having the CI and CD toolset sitting within a kubernetes cluster, rather than on puppetised and terraformed EC2 instances, it will allow for easier updates and upgrades to occur, and reduce the need for extra technology to handle configuration management (Puppet, Chef or Ansible) and infrastructure as code tools (Terraform)
- By having a single CI and CD toolset used by all tenants, it will enable Reliability Engineering to centrally manage and upgrade the toolset, rather than having tenants own the management and maintenance responsibility
