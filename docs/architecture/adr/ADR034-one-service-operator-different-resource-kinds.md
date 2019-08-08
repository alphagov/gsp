# ADR034: Single service operator for many different resource kinds

## Status

Accepted

## Context

In [ADR031](ADR034-one-service-operator-different-resource-kinds.md) we have
decided that we will be building an service operator for postgres.

We have came to an impasse in the discussion over whether we build separate
service operators for separate services - i.e. a single service operator for
Postgres, a single one for queues, etc. - or whether we should build a single
service operator that provisions access to all the services we need.

## Options

### Single service operator

Pros:
* Would only be one item on our current release pipeline
* Less repeating ourselves
* Would only need occupy one slot in reserved pods per node
* Easier to split out later
* Prevents divergance of operators accross the board
* Easier to add the next thing
* Monorepo all the things

Cons:
* All our CRDs would be parsed by something that is all-powerful
* Potential for it to become "THE" operator
* Could be solved by AWS soon - starting life at least as 'The missing parts of
  AWS Service Operator', with plans to later add e.g. gsp-local stubs

### Multiple service operators

Pros:
* Separation of concerns
* More idiomatic for Kubernetes as they're like microservices

Cons:
* Multiple new items in our release pipeline could exacerbate problems with
  Multi-tenant Concourse
* Would repeat ourselves a lot
* Each node has a limit to the number of pods - separate service operators
  would each need their own pod, taking up more resources and potentially more
  money
* Difficult to merge
* Risk of developing different incompatible operators as things diverge and new
  features are added to some but not others
* We would be making a network boundary where none is necessary

## Decision

We will make a single service operator for multiple services, and will be able
to split out parts as necessary if the considerations change over time. We
should keep an eye on its security, its scope, and AWS Service Operator
upstream activity.

