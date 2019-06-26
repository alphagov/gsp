# GDS Supported Platform Infrastructure

![overview of the GDS Supported Platform infrastructure](diagrams/gsp-architecture-infrastructure-1.svg)

<!--
__[edit draw.io diagram](https://www.draw.io/?state=%7B%22ids%22:%5B%221hUinA_Bejb-x9AGgso1iaBighXrCsIhJ%22%5D,%22action%22:%22open%22,%22userId%22:%22104206899246339571570%22%7D#G1hUinA_Bejb-x9AGgso1iaBighXrCsIhJ)__
-->


1. A GSP cluster resides in an AWS account within the London region (__eu-west-2__)

2. The infrastructure is deployed across three availability zones (__eu-west-2a__, __eu-west-2b__, __eu-west-2c__)

3. The cluster makes use of [Global Accellerator](https://aws.amazon.com/global-accelerator/) to manage static IP addresses.

4. Load balancing is achieved through the use of an [application load balancer](https://aws.amazon.com/elasticloadbalancing/features/#Details_for_Elastic_Load_Balancing_Products).

5. External egress access is via a [NAT gateway](https://docs.aws.amazon.com/vpc/latest/userguide/vpc-nat-gateway.html) in each availability zone.

6. The [kubernetes control plane](https://kubernetes.io/docs/concepts/#kubernetes-control-plane) is managed by [AWS EKS](https://aws.amazon.com/eks)

7. The cluster relies on [AWS IAM](https://aws.amazon.com/iam) for identity and authorisation

8. The cluster uses [kiam](https://github.com/uswitch/kiam) to allow cluster users to associate [AWS IAM](https://aws.amazon.com/iam) roles to kubernetes pods

9. The cluster contains an autoscaling group containing a continuous integration service based on [ConcourseCI](http://concourse.ci/)

10. By default there are three [kubernetes worker nodes](https://kubernetes.io/docs/concepts/architecture/nodes/) for the cluster split across the three availability zones

11. All accounts benefit from [AWS Shield](https://aws.amazon.com/shield/) protection against distributed denial of service attacks

12. The cluster aggregates selected log events to [AWS CloudWatch](https://aws.amazon.com/cloudwatch/) for ongoing processing by the Cyber Security team

13. CloudWatch logs are shipped externally to Splunk using [AWS Lambda](https://aws.amazon.com/lambda/)
