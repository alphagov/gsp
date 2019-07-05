# ADR021: Alerting

## Status

Accepted

## Context

The teams need timely notifications based on key indicators in order that they can ensure reliability and respond to issues.

The prometheus operator included in the GSP cluster can provide Alertmanager however we would like to manage alert routing across GDS and not duplicate routing rules or manage multiple sets of alert targets.

## Decision

We will route alerts to a separately hosted shared [Alertmanager](https://prometheus.io/docs/alerting/alertmanager/) to handle platform alert routing

## Consequences

- Using an external alertmanager will require additional configuration in each cluster
- We will be unable to take advantage of the [automated configuration](https://coreos.com/operators/prometheus/docs/latest/user-guides/alerting.html) and custom resources from the prometheus operator
- We will be able to use Pagerduty
