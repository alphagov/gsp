# ADR 003: Placement of CI and CD Tools

## Status

Proposed

## Context

The placement of the CI and CD toolset, either within or external to the control cluster and / or tenant cluster, determines most aspects of the build and deployment toolset and influences architectural decisions.


## Decision

- The CI and CD tools will run seperate from the control cluster
- The CI and CD tools will run within their own kubernetes cluster
- There will be a single CI and CD toolset used by all tenants of the new service


## Consequences

- With a single CI and CD toolset used by all tenants, there will be a need for those tools to have strong RBAC to prevent cross-tenant pollution
- The decision to use a single tool for all tenants could easily be adapted in the future to allow per-tenant CI and CD tool deployments should requirements change
- Authentication for the CI and CD toolsets should align with [ADR 002](https://github.com/alphagov/gsp-team-manual/blob/master/adr/002-identity-provider.md)
- It is believed that by having the CI and CD toolset sitting within a kubernetes cluster, rather than on puppetised and terraformed EC2 instances, it will allow for easier updates and upgrades to occur, and reduce the need for extra technology to handle configuration management (puppet, chef or ansible) and infrastructure as code tools (terraform)
- By having a single CI and CD toolsets used by all tenants, it will enable to RE to centrally manage and upgrade the toolset, rather than having tenants own the management and maintenance responsibility
