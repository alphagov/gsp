# ADR 005: Ingress

## Status

Accepted

## Context

Creating a [`Service`](https://kubernetes.io/docs/concepts/services-networking/service) resource of type [`LoadBalancer`](https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer) will result in cloud provider specific infrastructure – for example an [Elastic Load Balancer](https://aws.amazon.com/elasticloadbalancing/) on AWS – being provisioned to expose the Service to the public.

This is unfortunately a bit of an anti-pattern as:

* Provisioning one load balancer per exposed service would be unnecasarily expensive.
* Provisioning a load balancer takes time, which can result in significantly slower deployment times.
* LoadBalancer types of Service require vendor specific extensions and have often have vendor specific configuration that is not portable

## Decision

* We shall provide a single load-balancer to provide ingress to each cluster.
* We shall provide an [ingress controller](https://kubernetes.io/docs/concepts/services-networking/ingress/#ingress-controllers) as part of our managed cluster deployments to enable Ingress resources
* We shall provide guidance to promote the use of Ingress over [`LoadBalancer`](https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer) types of [`Service`](https://kubernetes.io/docs/concepts/services-networking/service) for routing
* We shall not support and may disable the use of [`LoadBalancer`](https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer) types of [`Service`] entirely

## Consequences

* The LoadBalancer Service type is well documented around the web and does provide a convinent way to expose a Service that people may expect to be enabled
