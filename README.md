
# What is the GDS Supported Platform?

## Overview

The GDS Supported Platform is the [strategic container platform](https://reliability-engineering.cloudapps.digital/documentation/strategy-and-principles/re-strategy.html#2-build-a-new-common-hosting-platform-for-gds-services) for hosting services developed for the [Government Digital Service](https://www.gov.uk/government/organisations/government-digital-service).

The platform manages the infrastructure that your application runs on and provides tooling for teams to build, deploy and manage their applications on that infrastructure.


The platform is operated by [GDS Reliability Engineering](https://reliability-engineering.cloudapps.digital/) according to the [Technology & Operations Shared Responsibility Model](https://reliability-engineering.cloudapps.digital/documentation/strategy-and-principles/techops-shared-responsibility-model.html)

## Who is the platform for?

The platform is for teams working in the [Government Digital Service](https://www.gov.uk/government/organisations/government-digital-service) that need to run software applications.

## Features

- A declarative continuous delivery workflow - merging to master triggers deployment to production
- A container platform based on industry standard [Docker](https://docs.docker.com/) and [Kubernetes](https://kubernetes.io)
- Build and release automation with [ConcourseCI](https://concourse-ci.org/)
- A private container registry with [Docker Registry](https://docs.docker.com/registry/)
- Signing of docker image integrity with [Docker Notary](https://docs.docker.com/notary/)
- Scanning of docker images for security vulnerabilities with [clair](https://github.com/coreos/clair)
- Monitoring and alerting with [Prometheus](https://prometheus.io/), [Alertmanager](https://prometheus.io/docs/alerting/alertmanager/) and [Grafana](https://grafana.com/)
- Secure git-based secrets configuration with [sealed-secrets](https://github.com/bitnami-labs/sealed-secrets)
- Ingress management with [nginx ingress controller](https://kubernetes.github.io/ingress-nginx/)
- Protective monitoring provided by GDS TechOps CyberSecurity with [Splunk](https://www.splunk.com/)
- Cloud infrastructure hosted on [AWS](https://aws.amazom.com) in three availability zones in the London region managed with [Terraform](https://www.terraform.io/)
- Kubernetes control plane with [AWS EKS](https://aws.amazon.com/eks/)

## Getting started
- [Host a GDS Supported Platform cluster in the AWS cloud](docs/gds-supported-platform/getting-started-gsp-cluster.md)
- [Host a GDS Supported Platform cluster locally](docs/gds-supported-platform/getting-started-gsp-local.md)



## Help and support
For help or support:
- read our [documentation](/docs)
- raise an [issue](https://github.com/alphagov/gsp-terraform-ignition/issues)
- message the team on the Reliability [Engineering Slack channel](https://gds.slack.com/messages/CAD6NP598) [#reliability-eng](https://gds.slack.com/messages/CAD6NP598)
