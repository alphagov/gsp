# ADR025: Ingress

## Status

Accepted

## Context

We currently have two [ingress][Ingress] systems:

* Istio (see [ADR019])
* nginx-ingress (see the old Ingress [ADR005])

Istio's [Virtual Service] records are essentially advanced `Ingress` records.

Do we need both?

## Decision

No. We will use an [Istio Ingress Gateway](https://istio.io/docs/tasks/traffic-management/ingress/ingress-control/)

## Consequences

* Less to manage
* [Ingress] is one of the standard kubernetes types, as such people will expect it to work

[ADR005]: ADR05-ingress.md
[ADR019]: ADR019-service-mesh.md
[Ingress]: https://kubernetes.io/docs/concepts/services-networking/ingress/
[Virtual Service]: (https://istio.io/docs/reference/config/networking/v1alpha3/virtual-service/
