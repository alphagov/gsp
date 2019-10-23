# ADR040: Ensuring cluster stability while replacing nodes

## Status

Accepted

## Context

SREs sometimes need to change various things about a cluster's worker nodes, and therefore replace all of the nodes:
* AMI (e.g. for an EKS update)
* instance type
* anything else in the launch template in some way (e.g. instance role)

In addition to this, in our Terraform we get the latest AMI ID for our EKS version from AWS - this means that as soon as AWS releases a new AMI for our current EKS version, the next time the cluster's deployer step runs for Terraform it will also replace all the nodes.

We have a CloudFormation UpdatePolicy on our worker node auto-scaling-groups that tells it to replace 1 node at a time in each ASG. This would be fine if we didn't have 3 ASGs with a small number of nodes in each - right now this practice terminates too many instances at once and results in the cluster becoming unstable. This in turn causes outages in the applications running in the cluster.

Without such a policy we'd be able to update the launch template for the ASG and it wouldn't remove any existing nodes, just set up any new future nodes correctly.

The reason we have 3 small ASGs is so the cluster auto-scaler can scale nodes independently in each AZ - it will know that a pod attempting (and failing) to schedule is tied to a persistent volume claim in a particular availability zone and can scale up the ASG for that zone.

## Possible actions

Here are some of the actions we may want to consider. It is expected we'll choose some subset/variation of these.

### Action 1: Make CloudFormation UpdatePolicy wait for new nodes to be healthy when scaling up

Specify WaitOnResourceSignals, MinSuccessfulInstancesPercent and PauseTime in the CloudFormation UpdatePolicy. Mutually exclusive with action 2. Likely also requires action 4 to be safe.

The [AWS docs](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-attribute-updatepolicy.html) say this about the WaitOnResourceSignals parameter:
```
Specifies whether the Auto Scaling group waits on signals from new instances during an update. Use this property to ensure that instances have completed installing and configuring applications before the Auto Scaling group update proceeds. AWS CloudFormation suspends the update of an Auto Scaling group after new EC2 instances are launched into the group. AWS CloudFormation must receive a signal from each new instance within the specified PauseTime before continuing the update. To signal the Auto Scaling group, use the cfn-signal helper script or SignalResource API.
```

Sometimes new worker nodes don't appear ready in Kubernetes until after the previous worker node they're replacing has been terminated. This should allow us to resolve that and ensure that for every instance terminated e.g. through launch template changes, we have a replacement already in-place.

We could use one of these approaches:
* When a worker node starts we could have it schedule a Job, selecting itself to run on, that calls SignalResource.
* We could have a DaemonSet that calls SignalResource and then sleeps infinitely

Note that SignalResource callers will need EC2 instance ID as well as the ID of the CloudFormation stack that created the ASG it comes from.

Pros:
* AWS should mostly take care of the process for us, we'd just need some code to call SignalResource.
* The cluster should not encounter any health issues if new instances are ready to accept new pods before existing instances get terminated

Cons:
* Our Concourse job would probably time out while this process happens, but we can probably adjust/live with that.

### Action 2: Custom process for rolling nodes instead of CloudFormation UpdatePolicy

Remove the CloudFormation UpdatePolicy and find some alternative to the update policy that ensures the cluster eventually fully rolls out the expected launch template. Mutually exclusive with action 1.

#### Action 2A: Manually roll nodes

Manual cordoning, draining and termination of nodes

Pros:
* Total control over each step in the process.

Cons:
* This would make it some rather irritating toil as an engineer would have to manually run the process.
* It might be possible to make a mistake here and cordon/drain too many nodes at once without termination so therefore no replacement.
* this requires `admin` privileges, which we'd like to avoid for routine maintenance tasks

#### Action 2B: Script for rolling nodes, triggered manually

A script that we run that cordons, drains and terminates nodes - triggered manually.

Pros:
* Would be able to choose exactly when we replace nodes to optimise for low load, no other critical operations ongoing, staff availability in the event of a problem, etc.

Cons:
* We'd have to notice the new launch template and trigger it manually.

#### Action 2C: Regularly terminate nodes

A script that runs regularly that terminates nodes - this would ensure the cluster is eventually updated, just not immediately.

Pros:
* Has the added benefit of regularly killing off nodes regardless of launch template changes, meaning nodes don't stick around for too long.

Cons:
* If this was a CronJob inside the cluster, it might terminate itself while running and leave us with an inconsistent cluster? But if it's a Lambda maybe not. Or if it's written to only kill one node and wait until the next run to kill another we might also be okay.
* Would not immediately replace all the nodes, rollout would be slow.
* If the new launch template is broken, the cluster would slowly break itself through terminating the remaining healthy nodes.

#### Action 2D: K8s node operator

Write an operator that manages (i.e. taking over control of the desired instances count of the ASG) the ASGs for the worker nodes that works with / replaces the cluster autoscaler. Mutually exclusive with action 3.

Pros:
* Would be able to manage both node replacement for a new launch template and scale for load at the same time.

Cons:
* More code for us to write and maintain.
* Potentially re-inventing the wheel if we have to replace the cluster autoscaler.

### Action 3: Use Auto Scaling lifecycle hooks

A lifecycle hook in the ASG that tells it to cordon and drain nodes etc. when it terminates them.

From the docs:

> [Auto scaling lifecycle hooks][lifecycle-hooks] enable you to perform custom actions by pausing instances as an Auto Scaling group launches or terminates them. When an instance is paused, it remains in a wait state until either you complete the lifecycle action using the complete-lifecycle-action CLI command or CompleteLifecycleAction API action, or the timeout period ends (one hour by default).

Essentially when a node is started, we can register some hooks so that when the ASG launches a node, it waits until the node is healthy in the cluster; and when it terminates a node, it waits until we've cordoned and drained it.

Pros:
* This might have some benefit if we want to e.g. use spot instances in the future, as we'd be confident we can quickly cordon and drain a node in the event of a notification saying a current worker node will be terminated shortly.
* This will work with any kind of autoscaling event - eg CloudFormation UpdatePolicy will just work with lifecycle hooks in place

Cons:
* None?

[lifecycle-hooks]: https://docs.aws.amazon.com/autoscaling/ec2/userguide/lifecycle-hooks.html

### Action 4: Avoid interference from cluster autoscaler

Scale autoscaler down so it's not running while we replace nodes. Mutually exclusive with action 2D.

This would likely need to be done every time Terraform apply is run (until that's run we don't know if nodes will be replaced), with checks to permit cluster creation and to avoid breaking if the cluster autoscaler deployment does not exist.
The pseudocode would look like this:
```
if eks-cluster exists:
    if cluster-autoscaler deployment exists:
        scale autoscaler deployment to 0 replicas

terraform apply

if cluster-autoscaler deployment exists:
    scale autoscaler deployment to 2 replicas
```

Pros:
* This is [AWS's recommended approach](https://docs.aws.amazon.com/eks/latest/userguide/update-stack.html).

Cons:
* Introduces some risk if there is a burst in pods wanting to schedule while we are replacing nodes.

## Decision

We will do action 3 (lifecycle hooks) (and therefore not action 2 (some other node rolling process)).  This will mean that the ASG will know how to wait for nodes to be ready in the cluster when launching, and how to wait for nodes to be drained before terminating.  It also means we can choose our method of node rolling with relative freedom.  We will keep our existing CloudFormation UpdatePolicy as our actual mechanism of rolling nodes.

We will do action 4 (turn off autoscaler while rolling nodes).  This will probably need to be a script run before we run terraform (to scale the autoscaler down) and afterwards (to scale the autoscaler back up).

## Consequences

* Node replacement should be much safer.
* Node replacement will likely be slower.
* Old instances should drain properly rather than suddenly disappearing.
* The cluster will not scale up in response to unschedulable pods while node replacement operations are in progress. This could cause problems if a deployment of lots of new pods occurs at the same time.
