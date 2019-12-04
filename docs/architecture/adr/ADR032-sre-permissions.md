# ADR032: SRE Permissions

## Status

Superseded by [ADR043](ADR043-k8s-resource-access.md)

## Context

As mitigation for some risks raised in threat modelling it was agreed that day-to-day access to the cluster was to be read-only for everyone. Only the concourse running in the cluster could make changes that originated from Github, which required several approvals before merging.

Following the gradual rollout of several applications onto the GSP it became clear there were issues with the deployment procedures. This caused conflicting and contending pods to attempt to execute, resulting in application failures and deployment pipeline blockages. This was happening up to several times a day, depending on the level of activity. The remedial procedure involves escalating one or more members to cluster admin to allow the resources to be deleted, before revoking the admin permissions again. This process requires 3 people to perform and could result in hours of wasted time for each occurrence.

## Decision

We will add to the SRE permissions map the ability to delete the following higher-level controllers so an escalation to cluster admin is no longer necessary:

* ConfigMap
* Deployment
* ReplicaSet
* Secret
* Service
* StatefulSet

We will also raise a story to investigate the root cause of the deployment issues with a view to removing these permissions in the future.

## Consequences

Deleting one of the above resources may result in downtime, depending on context, and will self-correct when the deployment pipeline for the application is run again.
