# ADR020: Metrics

## Status

Accepted

## Context

The teams looking after a cluster need visibility of key metrics in order that they can ensure reliability and diagnose issues.

[Prometheus](https://prometheus.io) and [Alertmanager](https://prometheus.io/docs/alerting/alertmanager/) are open-source systems monitoring and alerting tools and a graduated from the [Cloud Native Computing Foundation][CNCF].

A [kubernetes operator is available for Prometheus](https://github.com/coreos/prometheus-operator) that provides tight integration with the kubernetes API and minimal configuration required from service teams.

Reliability Engineering has standardised on [Prometheus](https://prometheus.io) to enable platform observability.

[CNCF]: https://www.cncf.io/announcement/2018/08/09/prometheus-graduates/

## Decision

We will use [Prometheus](https://prometheus.io/) and [Alertmanager](https://prometheus.io/docs/alerting/alertmanager/) managed by the [Prometheus Operator](https://github.com/coreos/prometheus-operator) for metrics in line with the standard [Reliability Engineering approach to Metrics and Alerting](https://reliability-engineering.cloudapps.digital/monitoring-alerts.html).

## Consequences

- Prometheus can use Alert Manager to generate alerts to notify engineers.
