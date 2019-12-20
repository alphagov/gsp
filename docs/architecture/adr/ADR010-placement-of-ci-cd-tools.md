# ADR010: Placement of CI and CD Tools

## Status

Superseded by [ADR029](ADR029-continuous-delivery-tools.md)

## Context

The placement of the CI and CD toolset, either within or external to the control cluster and / or tenant cluster, determines most aspects of the build and deployment toolset and influences architectural decisions.


## Decision

- The CI and CD tools will run separate from the control cluster
- The CI and CD tools will run within their own kubernetes cluster


## Consequences

- The separation of the CI and CD tool sets from the control cluster will ensure the control cluster remains as lean and allow us to make it as secure as possible.
