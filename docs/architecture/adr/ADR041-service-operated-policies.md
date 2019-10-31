# ADR041: Service operator provisioned policies 

## Status

Accepted

## Context

Our service-operator allows service teams to provision various AWS services by
declaratively defining resources and submitting them via the kubernetes api.

Some of these resources require IAM to authorise how the provisioned service
can be used. The types of actions that can be performed on.

#### Example

The service operator allows provisioning of S3 buckets and bucket configuration such as:

```
---
apiVersion: storage.govsvc.uk/v1beta1
kind: S3Bucket
metadata:
  name: s3-bucket-sample
spec:
  aws:
    LifecycleRules:
    - Expiration: 90days
    Versioning:
      Enabled: true
```

In order to access a provisioned bucket via the the AWS SDK users will require
an IAM role/policy that allows access.

We want things like bucket ACL, versioning configuration and lifecycle policy
to be defined declaratively via the resource manifest (see example above), and continuously managed
by the service operator.

We want users of the provisioned bucket to be able to read back all
configuration, and be able to fully utilise the specific bucket for reading,
writing and managing their objects within the provisioned bucket, but we want
to avoid giving permissions to users that could cause conflicts with the
properties that are managed by the service operator's reconcile loop. 

For example, given the example manifest above, we would like to avoid giving
permissions that would allow a user to alter the Expiration LifeCycleRules,
since any changes the user made would be periodically overwritten by the
service operator's reconciliation.

## Decision

* We will provision policy that gives full access for users to _use_ the
  provisioned service.
* We will avoid provisioning policy that allows users to create, destroy or
  configure the provisioned service, so that this can remain the declarative
  domain of the service-operator.

## Consequences

With only a single policy provisioned for each provisioned service, and no way
currently to request required permissions we may not be practicing the
"principle of least privilege".
