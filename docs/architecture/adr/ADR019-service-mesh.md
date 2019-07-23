# ADR019: Service Mesh

## Status

Accepted

## Context

Verify have a need to restrict exfiltration of data, enforce strict authentication between microservices and to use mutual TLS.

A service mesh gives us a way of meeting these needs.

### Option 1: Istio

Istio makes it easy to create a network of deployed services with load balancing, service-to-service authentication, monitoring, and more, with few or no code changes in service code.

Istio supports services by deploying a special sidecar proxy throughout your environment that intercepts all network communication between microservices, you then configure and manage Istio using its control plane functionality, which includes:

- Automatic load balancing for HTTP, gRPC, WebSocket, and TCP traffic.
- Fine-grained control of traffic behaviour with rich routing rules, retries, fail-overs, and fault injection.
- A pluggable policy layer and configuration API supporting access controls, rate limits and quotas.
- Automatic metrics, logs, and traces for all traffic within a cluster, including cluster ingress and egress.
- Secure service-to-service communication in a cluster with strong identity-based authentication and authorisation.

Pros/cons:

- an emerging standard (installed by default on GKE)
- a large community of contributors


### Option 2: AWS App Mesh (Istio from AWS)

[AWS App Mesh](https://aws.amazon.com/app-mesh/) is a service mesh that provides application-level networking to make it easy for your services to communicate with each other across multiple types of compute infrastructure. App Mesh standardizes how your services communicate, giving you end-to-end visibility and ensuring high-availability for your applications

pros/cons:

- Unavailable in London region
- Did not support automatic sidecar injection (meaning service teams would have to add lots of extra configuration to their Deployments)
- Appears to be abstraction over Istio


### Option 3: Linkerd 1.x & 2.0

[Linkerd](https://linkerd.io/) is an ultra light service mesh for Kubernetes. It gives you observability, reliability, and security without requiring any code changes.

Pros/cons:

- 1.0 has a richer feature set but poorer kubernetes support
- 2.0 has a very minimal feature set but native kubernetes support
- Going through major rewrite for improved Kubernetes support
- Smaller community
- Fewer features around


## Decision

We will use [Istio](https://istio.io/) to provide a service mesh in the GDS Supported Platform.

## Consequences

- Increased resource demands on nodes owing to use of sidecar containers
- Potential increase in the complexity of application configuration manifests
- Potential duplication of services provided by existing GSP components (i.e. Ingress controller)
