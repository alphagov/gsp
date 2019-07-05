# ADR018: GSP Local Development Environment

## Status

Accepted

## Context

Teams using the GDS Supported platform require the ability to develop, test applications and prove conformance with the GDS Supported platform on local hardware. Teams need to learn how to use the GSP and to understand how applications are containerised, packaged and deployed to a cluster using the standard CICD tools provided by GSP.

## Decision

We will [provide a way to run a full GSP compatible stack locally on a developer machine](/docs/gds-supported-platform/getting-started-gsp-local.md) without the cloud provider specific configuration.

## Consequences

- Lack of local machine resources (RAM and CPU) may be an issue due to deploying the full GSP stack
- Docker performance may slow down development and will require some effort to optimise
- The current environment lacks higher level tooling to streamline the workflow of the developer
