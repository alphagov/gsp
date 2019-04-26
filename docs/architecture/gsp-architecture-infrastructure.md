# GDS Supported Platform Infrastructure

![overview of the GDS Supported Platform infrastructure](diagrams/gsp-architecture-infrastructure-1.svg)

<!--
__[edit draw.io diagram](https://www.draw.io/?state=%7B%22ids%22:%5B%221hUinA_Bejb-x9AGgso1iaBighXrCsIhJ%22%5D,%22action%22:%22open%22,%22userId%22:%22104206899246339571570%22%7D#G1hUinA_Bejb-x9AGgso1iaBighXrCsIhJ)__
-->


1. The GSP cluster resides in a single AWS account within the London region (__eu-west-2__)

2. The infrastructure is deployed across three availability zones (__eu-west-2a__, __eu-west-2b__, __eu-west-2c__)

3. The [kubernetes control plane](https://kubernetes.io/docs/concepts/#kubernetes-control-plane) is managed by [AWS EKS](https://aws.amazon.com/eks)

4. By default there are three worker nodes for the cluster split across the three availability zones

5. The cluster relies on [AWS IAM](https://aws.amazon.com/iam) for identity and authorisation

6. The cluster uses [kiam](https://github.com/uswitch/kiam) to allow cluster users to associate [AWS IAM](https://aws.amazon.com/iam) roles to kubernetes pods

6. All accounts benefit from [AWS Shield](https://aws.amazon.com/shield/) protection against distributed denial of service attacks

7. The cluster aggregates log events to CloudWatch for ongoing processing by the Cyber Security team
