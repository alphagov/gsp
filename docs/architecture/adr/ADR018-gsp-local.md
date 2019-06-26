# ADR018: GSP Local Development Environment

## Status

Pending

## Context

Teams using the GDS Supported platform require the ability to develop, test applications and prove conformance with the GDS Supported platform on local hardware. Teams need to learn how to use the GSP and to understand how applications are containerised, packaged and deployed to a cluster using the standard CICD tools provided by GSP.

## Decision

We will provide a way to run a full GSP compatible stack locally on a standard GDS developer OSX based laptop based on a single node [minikube](https://github.com/kubernetes/minikube) cluster. 

## Consequences
- Developers using linux and windows may encounter platform specific issues due to lack of cross platform testing
- Lack of local machine resources (RAM and CPU) may constrain the use of the environments
- Docker performance on OSX may slow develpment
- The current environment lacks higher level tooling to streamline the workflow of the developer
