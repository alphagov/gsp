# ADR022: Logging

## Status

Accepted

## Context

We have various log sources:

* The EKS control plane (audit logs, api service, scheduler, control-manager)
* VPC Flow logs
* Load Balancer
* Pod/Application logs
* CloudHSM

All of these with the exception of the Pod/Application logs are stored by AWS in [CloudWatch](https://aws.amazon.com/cloudwatch/).

We would like a single storage location for indexing and search our logs for auditing and debugging purposes.

GDS currently have several common storage locations for logs:

* Logit.io (a SaaS ELK stack provider)
* Self hosted ELK stacks
* CloudWatch
* Splunk

Options:

### Option 1:

We could ship the Cloudwatch logs to logit.io using AWS lambda and ship the Pod/Application logs to Logit.io using something like [fluentd](https://www.fluentd.org/). This would assume that all users of the platform have a Logit.io instance and would end up duplicating a large number of the logs in both CloudWatch and Logit.io

### Option 2:

We could host a dedicate ELK stack (either in cluster or from AWS's managed offering) and ingest logs from both Pods and CloudWatch into the ELK stack. Managing ELK stacks has been a maintenance burden at GDS previously and this would require duplicating logs already stored in CloudWatch.

### Option 3:

We could ship the Pod/Application logs to CloudWatch using [fluentd](https://www.fluentd.org/) and expose CloudWatch insights interface to users of the platform

### Option 4:

We could ship the CloudWatch logs to Splunk using AWS lambda and ship the Pod/Application logs to Splunk using something like [fluentd](https://www.fluentd.org/). This would assume that all users of the platform have a Splunk instance and would end up duplicating a large number of the logs in both CloudWatch and Splunk.

## Decision

We will use [fluentd](https://www.fluentd.org/) to ship pod/application logs to [AWS CloudWatch](https://aws.amazon.com/cloudwatch/) to aggregate all platform/application logs to avoid double spending on log storage.

## Consequences

- some people like Kibana
- Cloudwatch insights may not meet the needs of day-to-day application debugging.
