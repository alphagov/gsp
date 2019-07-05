# ADR002: Containers

## Status

Accepted

## Context

At the time of writing the infrastructure/deployment landscape is:

* Many service teams are deploying applications to Virtual Machines (AWS EC2, VMWare, etc)
* Some service teams are deploying applications as containers (AWS ECS, GOV.UK PaaS, Docker)
* Few service teams are deploying applications as functions (AWS Lambda)

There is a mix of target infrastructure/providers in use, but there is a gradual migration towards hosting on AWS.

## Decision

We will focus on providing the primitives to run stateless containerised workloads.

## Consequences

* Some applications that were previously being deployed to Virtual Machines may require significant modification before they are suitable for running within a container.
* Some of the isolation guarantees present with Virtual Machines are not available to containers, which may limit our options when thinking about multi-tenancy or multi-environment architectures.
* An initial lack of solution to deploying event-based or Function-as-a-Service architectures may not be attractive to teams who are already experimenting in this area (however such systems _could_ be implemented on top of a container-based system)
