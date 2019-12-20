# ADR030: AWS Service Operator for SQS queues

## Status

Superseded by [ADR031](ADR031-postgres.md) and [ADR034](ADR034-one-service-operator-different-resource-kinds.md)

## Context

[Amazon announced](https://aws.amazon.com/blogs/opensource/aws-service-operator-kubernetes-available/) an AWS Service Operator for Kubernetes in October 2018.
There is a need for one or more users to provision SQS queues for their application to interact with.
The AWS Service Operator consists of a container that sits in the Kubernetes cluster and monitors for custom resource types which it maps to the appropriate CloudFormation resources and deploys.
It supports the following resources:
* CloudFormation Template
* DynamoDB
* ECR Repository
* ElastiCache
* S3 Bucket
* SNS Subscription
* SNS Topic
* SQS Queue

## Alternative ways of provisioning SQS

We considered some alternative ways of provisioning SQS queues. In particular, there is also an [AWS service broker](https://github.com/awslabs/aws-servicebroker) which provides AWS services via an [open service broker](https://www.openservicebrokerapi.org/) API. The use of service brokers is also interesting with a view towards future work.  In particular, GOV.UK PaaS provides services via service brokers, and so there would be opportunity for reuse of components such as [paas-rds-broker](https://github.com/alphagov/paas-rds-broker).

However, the current way that you integrate the service broker API into Kubernetes is via the [service catalog](https://github.com/kubernetes-sigs/service-catalog).  This is a somewhat leaky abstraction in that it adds a whole bunch of extra complicated Custom Resource Definitions (CRDs) that users have to understand; it also introduces a separate command-line interface for browsing the service catalog, which increases the number of tools a user has to learn.

We also considered that paas-rds-broker is currently deployed using bosh, which is a fairly invasive tool that we don't want to adopt if we can avoid it. However, after speaking with the GOV.UK PaaS team, they think that it would be relatively easy to dockerise paas-rds-broker and deploy it to kubernetes; so this is not a concern for our decision.

## Decision

We have included AWS Service Operator as part of GSP so that users can provision SQS queues.

## Consequences

* Users will be able to provision SQS queues using kubeyaml, without requiring support from RE
* Users will not have a way of provisioning queues in gsp-local; this is a recognised limitation for expediency right now but we need to revisit the local development experience
