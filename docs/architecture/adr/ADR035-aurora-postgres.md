# ADR035: RDS Aurora Postgres

## Status

Accepted

## Context

Our service operator will be responsible for providing Postgres along any other
services to our end users.

In this case, the ask is specifically Postgres, meaning we can pick our
solution for providing the instance as long as it exposes the correct APIs.

A reasonable set of candidates are:

- AWS RDS Postgres
- AWS RDS Aurora Postgres

### AWS RDS Postgres

A solution most comonly used across the board for databases. Is managed and
maintained by AWS, provides Backups and Snapshots.

### AWS RDS Aurora Postgres

Has the benefits of AWS RDS Postgres, and some more benefits. Such as:

- Scalable persistance storage
- IAM authentication
- Automatic Failover with multi AZs
- Faster read loads
- Slightly cheaper

Also has few downsides:

- Not recommended for heavy write systems (Twitter big)
- Slightly behind version wise

## Decision

We will continue with Aurora as we we don't have any specific requirements not
to and can benefit from the solution.

## Consequences

We may find ourselves in need of adding AWS RDS Postgres anyway in a future.
