# ADR031: Providing Postgres to tenants

## Status

Accepted, Supplemented by [ADR034](ADR034-one-service-operator-different-resource-kinds.md)

## Context

We have GSP tenants who want to deploy postgres instances for their
apps to use.

We have a techops design principle that we should use cloud-provided
stateful services where possible, rather than rolling our own.  This
means we would like to deploy postgres on RDS rather than running it
within kubernetes.

We have considered a number of options for how we might do this:

### Option 1: Terraform snowflaking

In this option, we write some custom terraform which gets applied in a
particular GSP cluster to provision specific RDS instances.  We then
need some configuration to make the credentials for the particular RDS
instance available to the application.

This option has the benefit of expediency (it would be quick to
implement) but it means that each time a tenant wants a new database,
they have to modify the terraform at the cluster level, so it doesn't
scale well.  It also is against the grain of kubernetes ways of doing
things.

### Option 2: AWS Service Operator

AWS have written
[aws-service-operator](https://github.com/awslabs/aws-service-operator),
a kubernetes operator which provides Custom Resource Definitions
(CRDs) for various AWS types.  It provides native types for various
things like SQS.  We have already decided to use it for SQS in
[ADR030](ADR030-aws-service-operator.md).

However, the AWS Service Operator doesn't support RDS.  There is an
open issue for RDS (awslabs/aws-service-operator#39) but activity is
slow.  This means we need to find another way to use the service
operator to provision RDS.

#### Option 2a: Use AWS Service Operator to deploy custom CloudFormation

The aws-service-operator provides a CRD to deploy a custom
CloudFormation template.  However, after we did some digging into
this, we realised that it *doesn't* provide a way of instantiating
that template into a stack.  It's hard to see how this is a useful
feature on its own.

#### Option 2b: Extend AWS Service Operator with RDS resource type

This would mean making a fork and adding our own RDS CRD to it.

We suspect that if we went down this road, we wouldn't be able to
contribute it back upstream, because the efforts in the
aws-service-operator seem to be focused on automatically generating
code for all resources from the CloudFormation Resource Specification
(see awslabs/aws-service-operator#153).  Our custom code would not fit
in with this design vision.

### Option 3: Build our own postgres operator

With this option, we would create our own Postgres CRD and build an
operator (using [kubebuilder](https://book.kubebuilder.io/)) to watch
for Postgres resources and provision RDS instances based on those
resources.

We choose a Postgres CRD, rather than an RDS CRD, because:

 - it's closer to the user's language ("I want a postgres" not "I want
   an RDS")
 - it's (in principle) more portable (other platforms provide
   postgres)
 - we have no need for other features of RDS right now (such as MySQL)

Within this, we have options for how we provision the RDS instance:

 - we could provision it with terraform
 - we could provision it with CloudFormation
 - we could provision it with raw AWS SDK code
 
If we use terraform, we need to find somewhere for the state file to
live (such as an S3 bucket).

If we use raw AWS SDK code, we have to implement our own logic for
resource updates (such as if a tenant wants to make their database
bigger).

CloudFormation offers the best of both: it has logic for resource
updates but no need for us to manage a statefile somewhere.

### Option 4: service brokers

There are service brokers available that we could use.  These are apps
built on the [open service broker
API](https://www.openservicebrokerapi.org/) which provide an interface
for provisioning postgres on RDS.  In order to use a service broker
with kubernetes, you need a bridging layer such as the [service
catalog](https://kubernetes.io/docs/concepts/extend-kubernetes/service-catalog/).

Two service brokers that we could use are the [aws service
broker](https://github.com/awslabs/aws-servicebroker) and our own
[paas-rds-broker](https://github.com/alphagov/paas-rds-broker).

It seems like a good thing to be able to reuse these service brokers.
We already deploy and run service brokers as part of GOV.UK PaaS, so
there is prior art and experience in the organisation for them.

However, the Service Catalog is a leaky abstraction and doesn't really
fit well in the kubernetes ecosystem as far as we can see.  It seems
geared towards an environment where you have hundreds of service
brokers, and provides separate CLI tools just for finding service
brokers and talking to them (rather than keeping everything within the
Kubernetes API and `kubectl`).  It also adds lots and lots of CRDs for
dealing with service brokers, rather than just the Postgres CRD that
other options might offer.

### Other context: Portable Service Definitions

There is a proposal within kubernetes to offer Portable Service
Definitions (see kubernetes/enhancements#706).  This would define
custom types at a higher level that different environments could
implement in different ways.  This provides interesting prior art in
how we might structure a Postgres CRD of our own.  We would also be
open to adopting the Postgresql Standard Resource Definition (SRD)
(see libresh/StandardResourceDefinitions#6) in future, but not just
now.  We anticipate that we can build something focused on our
immediate needs now, but that it would not be terribly expensive to
adopt the SRD in future.

## Decision

We will go with Option 3 above.  That is, we will create a custom
postgres-operator and Postgres CRD which looks for Postgres resources
and provisions RDS instances based on this.  The operator will output
the credentials for accessing the RDS instance as a Secret in the
requesting namespace.

## Consequences

Users will be able to provision instances of postgres just by creating
a Postgres kubeyaml object in their namespace.  In particular, if they
want to deploy their app to a new namespace, they will be able to
without RE intervention.

Access to Postgres will be controlled by posession of the username and
password.  This will be controlled by putting the username and
password in a Secret in the given namespace.  Tenants in different
namespaces will therefore not be able to access each others' postgres
instances.

This may take longer than we anticipate, in which case we have a
fallback option of Option 1 above for expediency.
