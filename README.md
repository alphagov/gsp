# GSP [![IRC](https://img.shields.io/badge/kubernetes-v1.16-0099ef.svg)]() <img align="right" src="./docs/assets/gsp.png" alt="gsp" width="30%" height="whatever">

***
This project solved some specific needs of GDS. It was not generally useful for people outside of GDS. You should consider using [GOV.UK PaaS](https://www.cloud.service.gov.uk/) if you are looking for somewhere to run your services.
This is a decommissioning notice detailing issues that would need to be solved in order to re-use this codebase. It only documents issues known at time of repository archiving, when it will cease being updated.
For the old README prior to archiving, see [README-old.md](/README-old.md)
***

GSP ([GDS](https://www.gov.uk/government/organisations/government-digital-service) Supported Platform) was a Kubernetes distribution based on Amazon EKS.

Technically:
* The Kubernetes/EKS version is behind - we're on 1.16 and the latest is 1.20. 1.16 will not be possible to use after July 2021.
 * There's a TODO in `pipelines/deployer/deployer.yaml` about k8s 1.15 we can probably remove.
* GSP relies on Istio 1.5.8 which became end-of-life on 2020-08-24. Also 1.6 was end-of-life on 2020-11-23, 1.7 was end-of-life on 2021-02-25, and 1.8 will be end-of-life on 2021-05-12.
* GSP ran Prometheus and Grafana through prometheus-operator 8.15.6, and it's no longer developed at `git@github.com:helm/charts.git` stable/prometheus-operator - the latest version of chart we were using is now deprecated.
* The check-vulnerabilities job in each cluster deployment pipeline would find all sorts of things in the third-party images we used like cluster-autoscaler, concourse-web, external-dns, fluentd-cloudwatch, fluentd-kubernetes-daemonset, and more. Some of these may be resolvable by upgrading the version of the software used.
* It's based on Terraform but with some weird extra CloudFormation that should be merged into the Terraform - `modules/k8s-cluster/data/nodegroup-v2.yaml` and (especially obscure, but possibly unnecessary depending on the item below) `modules/k8s-cluster/data/nodegroup.yaml`.
* We're not completely certain that the cluster-management nodes still necessary, it may be possible to put gatekeeper and the cluster-autoscaler on normal worker nodes.
* There is a strange distinction between the k8s-cluster and gsp-cluster terraform modules which should probably be eliminated.
* Some of the docs are written with gds-cli in mind but gds-cli dropped support for GSP in version 4, so you'll either need to convince gds-cli to reimplement that support or eliminate the dependency by going back to plain aws-vault.
* We wrote our own service operator, but then AWS made https://aws.amazon.com/blogs/containers/aws-controllers-for-kubernetes-ack/ which probably obsoletes part of it already - and in future, likely all of it. If you're going to re-use the GSP code, remove the obsolete parts of our service operator (ECR, SQS, S3) and use AWS's, then eventually remove ours all together when they add RDS+ElastiCache
* We tried to have working, daily replacement of the underlying EC2 instances, but our experience with node rolling in GSP was complicated, with multiple incidents. There's some problems lurking somewhere that need tracking down and solving. It may be possible to use EKS Fargate and eliminate the EC2 instances.
* GSP did not support automatically getting certs for non-govsvc.uk domains (i.e., subdomains of service.gov.uk), leading us to do some odd tricks involving separately terraforming CloudFront distributions in front of GSP. You would probably want to solve this if you use it.
* The `lambda_splunk_forwarder` terraform module still lurks but new clusters should be using CSLS to ship logs, so this Lambda should be eliminated and not re-used. It only survived until GSP decommissioning because swapping it out on the existing production cluster may have been an issue that lead to dropped logs.
* There's still some CloudHSM support lurking in GSP, this should be eliminated and not re-used.
* It was set up for the govsvc.uk domain which may have to be re-reigstered or transferred.
* Permissions issues [documented here](https://docs.google.com/document/d/1za1H8XZaDd9LUUpOd0fmbwY3EOz35G39d7xJ5F_e6A4/edit?usp=sharing)

More broadly:
* Ultimately, an organisation needs strong justification to run multiple platforms like this. If we were to spin this back up it'd have to replace existing systems.
* Kubernetes is complicated from a developer point of view, requiring lots of tricky YAML documents.
* The security model was motivated by a particular high-security project which made it tricky to get things done as everything needed to go through GitOps.