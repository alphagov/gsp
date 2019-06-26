# ADR029: Pull based Continuous Delivery tools

## Status

Accepted

## Context

We currently have a separate GSP cluster for CI tooling that runs Concourse and other build tools (Docker Registry, Image Scanning and Notary)

We currently share configuration between the tools cluster and programme clusters.

We have in-cluster tooling (based on [Flux](https://github.com/weaveworks/flux)) to perform a pull based continuous deployment loop (watch a branch, apply state to cluster).

We forked flux to support git commit signature verification.

We want to avoid sharing credentials capable of modifying the cluster state outside of the cluster itself.

We have been experiencing the following issue with the current setup:

* Poor visibility to service teams of what Flux is doing
* Slow and hard to reason about issues with deployments
* Complex configuration that needs to span multiple cluster instances reduces our ability to automate provisioning/deployments
* Having separate credentials for the tools cluster and programme clusters is confusing (they both look the same) and inconvenient (different kubeconfig, different IAM profile etc.)

If we run the CI/CD tools inside the programme cluster it could:

* Remove the need to have separate CD tooling (use Concourse for both CI and CD tasks)
* Remove the need to maintain our own fork of Flux (we can achieve the same thing with the Concourse resource we already have)
* Improve visibility of continuous deployment by using Concourse's user interface to expose status to developers
* Reduce the amount of configuration required (no need to share config between clusters/accounts)
* Automate configuration of build tooling and provenance checking (signatures, Notary, etc)
* Reduce the amount of infrastructure we manage to one control plane, one account etc.
* Automate access control to concourse using the same RBAC rules as the cluster (One Concourse Team == One Namespace in the cluster)

## Decision

We will run Concourse in each programmes's cluster to provide both CI and CD tooling.

## Consequences

* We may need to isolate the worker nodes of the CI/CD from the regular applications worker as workloads are suitably different
* Not all CI tasks may be suitable to run in such a production environment (ie PR builds from untrusted sources)
