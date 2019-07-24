# ADR030: AWS Service Operator

## Status

Pending

## Context

[Amazon announced](https://aws.amazon.com/blogs/opensource/aws-service-operator-kubernetes-available/) an AWS Service Operator for Kubernetes in October 2018.
There is a need for one or more users to provision queues (potentially SQS ones) for their application to interact with.
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

## Decision

We have included AWS Service Operator as part of GSP so that DCS can use it.

## Consequences

* Applications may become dependent on it and struggle to run with gsp-local
* Potential for vendor lock-in?
