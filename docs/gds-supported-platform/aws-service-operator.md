# Using AWS Service Operator

## How to use it
AWS Service Operator is a tool we use to allow GSP users to write kubeyaml resources that will directly map to AWS resources (via CloudFormation that it generates and deploys). The end result of this is that developers can create certain types of AWS resources directly inside their kubeyaml. Here's an example:
```yaml

- kind: SQSQueue
  apiVersion: service-operator.aws/v1alpha1
  metadata:
    name: alexs-test-queue
  spec:
    messageRetentionPeriod: 3600
    maximumMessageSize: 1024
```
This will create an SQS Queue on AWS named alexs-test-queue, with a message retention period of 1 hour, and a maximum message size of 1KiB.
Note that missing those parameters from the spec may lead to it trying to deploy CloudFormation that simply errors as soon as AWS attempts to run it. If a resource does not appear as expected you can go into the AWS CloudFormation console and find out why.

Alongside SQS Queues, it supports the following resources:
* CloudFormation Template
* DynamoDB
* ECR Repository
* ElastiCache
* S3 Bucket
* SNS Subscription
* SNS Topic

## How it works
You don't need to know this to use it, this information is for cluster operators.
AWS Service Operator consists of a container that runs essentially a daemon, and a (slightly modified in our case) kubeyaml config that sets up the container, provides a bunch of custom resource definitions (e.g., there is a definition in there for SQS Queues), etc. - it also gives the container access to monitor stuff going on around the cluster.
The daemon monitors the k8s cluster for such custom resources being created and will deploy to AWS the relevant CloudFormation to create the requested resource.
