# Updating EKS

## Overview

Patch-level EKS version upgrades are performed automatically by AWS, and do not
require human intervention.

Major/Minor-level EKS version upgrades require us to explicitly set the version
in our deployment configuration, manually update some of the control-plane
components and ensure that worker nodes are running compatible versions.

AWS publishes guidance for upgrades [on the EKS documentation
site](https://docs.aws.amazon.com/eks/latest/userguide/update-cluster.html).

We try to split the upgrade process into discrete stages so that the impact of each stage can observed:

1. Upgrade the control-plane EKS version
1. Upgrade control-plane components (if required)
1. Upgrade worker nodes

All three stages should be performed and tested in an on-demand cluster before an upgrade to any of the components is merged to master and rolled. This ensures we don't release any time-bombs.

## Upgrade EKS version

The control plane version is set from a terraform variable.

The value for this variable comes either from the cluster's "cluster-config" or
from the default values provided by the gsp deployer pipeline.

Most clusters do not pin a specific EKS version, so updating the default
`eks-version` is usually all that is required to upgrade all clusters.

* See [PR568: upgrade EKS control plane](https://github.com/alphagov/gsp/pull/568/files) for an example of updating default control-plane version.

## Upgrade control plane components

As part of the [EKS upgrade documentation](https://docs.aws.amazon.com/eks/latest/userguide/update-cluster.html)
AWS _may_ recommend that some other components are upgraded to align with
control-plane changes.

We maintain these components as templates within the gsp-cluster chart.

* See [PR541: upgrade coredns](https://github.com/alphagov/gsp/pull/541/files) for an example of upgrading DNS components.
* See [PR534: upgrade vpc-cni and calico](https://github.com/alphagov/gsp/pull/534/files) for an example of upgrading networking components.
* See [PR542: upgrade kube-proxy](https://github.com/alphagov/gsp/pull/542/files) for an example of updating the KubeProxy components

## Upgrade worker nodes

We deploy worker nodes as EC2 instances based on the [Amazon EKS-Optimized
Linux AMI](https://github.com/awslabs/amazon-eks-ami) in an AutoScalingGroup
losely based on the cloudformation template provided by amazon, but deploy it
via terraform.

AWS will advertise the AMI ID that is compatible with a given EKS control-plane
version [in their documentation](https://docs.aws.amazon.com/eks/latest/userguide/eks-optimized-ami.html).

Changes to the cloudformation template (such as updating the AMI reference)
trigger a rolling deployment of each worker node on deploy.

An [AutoScaling lifecycle hook](https://github.com/alphagov/gsp/tree/master/components/aws-node-lifecycle-hook)
helps ensure that instances are drained before retirement as the rolling
deployment takes place.

* See [PR1134: Bump worker nodes AMI](https://github.com/alphagov/gsp/pull/1134/files) for an example of changing the worker node EKS version.

## Final touches

There's a badge on the README.md that notes the current kubernetes version. Give it a bump!
