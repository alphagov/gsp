# ADR024: Soft Multi-tenancy

## Status

Accepted

## Context

One Programme has many Service Teams.

One Service Team has many Environments.

Some Service Teams have separate AWS accounts for separate environments. (i.e. Staging, Production)

Many Service Teams have micro-service architectures that run on many machines.

Some Service Teams have unique programme specific network isolation requirements that may be hard to implement in a shared environment.

Separate programme level accounts would enable separation of billing.

Sharing the infrastructure within a programme will lower hosting costs.

To ensure network/compute isolation between Service Teams it may be necessary to isolate resources.


## Decision

We will design for a "soft multi-tenancy" model where each programme shares a single GSP cluster with service teams within that programme.

This will:

* Maintain clear separation of billing at the programme level by isolating cluster to programme's own AWS account
* Maintain clear separation of programme specific policies and risk assessments by not forcing all users to adhere to the strictest rules?
* Minimize costs by sharing infrastructure, control plane and tooling between teams/environments
* Minimize support burden by reducing the amount of configuration

## Consequences

* Less efficient than one big cluster
* Less isolated than millions of clusters
