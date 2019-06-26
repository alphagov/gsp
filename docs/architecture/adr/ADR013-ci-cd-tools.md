# ADR013: CI & CD Tool

## Status

Accepted

## Context

We need to choose which tool or tools to use for CI and CD.  Different tools suit different purposes, however some cross over exists which could allow the use of a single tool to do both CI and CD.

## Decision

We will use [Concourse](https://concourse-ci.org/) for both CI and CD.  

Reasons:

- It will allow the alpha work to progress without waiting for a decision based upon user research on which tool set is best suited for use for CI or CD
- The team has experience of using concourse for CI and CD with kubernetes
- A working example already exists that can be extended for use in the alpha
- Concourse supports simple RBAC which should allow for multi-tenancy capability in the future
- It will accelerate the development of the alpha, with the team only needing to learn a single tool rather than multiple tools


## Consequences

- Research to discern whether Concourse is still an appropriate choice will occur as team becomes better familiar with the technology
- User research will be possible once an alpha can be presented to users for evaluation, and a decision to change tool or begin using different tools for CI and CD will be possible
