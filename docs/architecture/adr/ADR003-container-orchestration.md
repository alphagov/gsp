# ADR003: Container Orchestration

## Status

Accepted

## Context

[ADR002](ADR002-containers.md) decided that we should provide the primitives for building, deploying and running container-based workloads.

In order to meet the needs existing service teams such a primitive must be able to support:

* **High Availability** and **Load Balancing**: running multiple instances of applications across multiple physical nodes
* **Horizontal Scalability**: ability to increase/decrease the number of instances of an application to respond to demand
* **Declarative Configuration**: ability to describe infrastructure/deployments as code
* **Internal Networking**: ability to restrict traffic to a non internet-facing network
* **Service Discovery**: ability to describe relationships between running containers as microservices

There are also some expectations that people may have from experience with other container-based platforms, for instance:

* **Self Healing**: features that aid with the restarting and replacing of failed application instances
* **Logging**: solutions to how to access, and ship log output from running applications

There are many popular solutions that can meet the needs of this orchestration problem:

* **Vendor specific**: Cloud Providers have various solutions to manage various container workloads. AWS ECS/Fargate is currently in use by a few service teams at the time of writing as could meet the needs.
* **CloudFoundry**: It provides a full-stack PaaS-type solution for building and running container workloads. GOV.UK PaaS is based on CloudFoundry.
* **Cloud Native** / **Kubernetes**: Provides a minimal but extensible API for managing container workloads and large community of compatible projects for common solutions.
* **OpenShift**: Builds upon Kubernetes to provide a full-stack PaaS-type solution for building and running container workloads.
* **DCOS/Marathon**: An extra abstraction layer (DCOS) above a container orchestrator layer (Marathon) that could potentially support other types of workloads ie VMs.
* **Docker Swarm**: An orchestrator tightly integrated with the Docker container runtime.

In [ADR001](ADR001-support-model.md) we decided against running a "batteries included" PaaS-type platform, and to sacrifice the centralised control such a system would bring in exchange for giving more control/flexibility to service teams.

We did not feel like the additional abstraction layer in DCOS/Marathon provided any value we needed.

We were unsure that a system based on Swarm would provide the flexibility we want to respond to the evolving needs of service teams.

Kubernetes has been emerging as an industry standard, with most major cloud providers offering various services around the technology. It offers the basic building blocks for managing container workloads but has a large strong community of compatible projects ([CNCF](https://landscape.cncf.io/)) to extend functionality.

The vendor specific solutions such as ECS/Fargate are very compelling, potentially reducing the maintenance burden of running infrastructure. However as a public body we should try to avoid unnecessary vendor lock-in.

## Decision

We will build upon Kubernetes as our base container orchestrator and IaaS abstraction layer.

## Consequences

* Kubernetes offers only the very low-level primitives, we will likely need to provide solutions to several problems that would have been bundled along with a more full-stack PaaS offering.
