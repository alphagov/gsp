# ADR030: AWS Service Operator for SQS queues

## Status

Accepted

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

## Decision

We have included AWS Service Operator as part of GSP so that users can provision SQS queues.

## Consequences

* Users will be able to provision SQS queues using kubeyaml, without requiring support from RE
* Users will not have a way of provisioning queues in gsp-local; this is a recognised limitation for expediency right now but we need to revisit the local development experience
