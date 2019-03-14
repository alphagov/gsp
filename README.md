# What is the GDS Supported Platform?

## Overview
The GDS Supported Platform is the strategic container platform for hosting services developed for the [Government Digital Service](https://www.gov.uk/government/organisations/government-digital-service).

The platform manages the infrastructure that your application runs on and provides tooling for teams to build, deploy and manage their applications on that infrastructure.


The platform is operated by [GDS Reliability Engineering](https://reliability-engineering.cloudapps.digital/) according to the [TechOps shared responsibility model](https://github.com/alphagov/gsp-team-manual/blob/master/docs/gsp-shared-responsiblity-model.md)

## Who is the platform for?

The platform is for teams working in the [Government Digital Service](https://www.gov.uk/government/organisations/government-digital-service) that need to run software applications.

## Features

- A declarative continuous delivery workflow - merging to master triggers deployment to production
- A container platform based on industry standard [Docker](https://docs.docker.com/) and [Kubernetes](https://kubernetes.io)
- Build tooling with [ConcourseCI](https://concourse-ci.org/)
- A private container registry with [Docker Registry](https://docs.docker.com/registry/)
- Signing of docker image integrity with [Docker Notary](https://docs.docker.com/notary/)
- Scanning of docker images for security vulnerabilities with [clair](https://github.com/coreos/clair)
- Application deployment with [Flux](https://github.com/weaveworks/flux)
- Monitoring and alerting with [Prometheus](https://prometheus.io/), [Alertmanager](https://prometheus.io/docs/alerting/alertmanager/) and [Grafana](https://grafana.com/)
- Secure git-based secrets configuration with [sealed-secrets](https://github.com/bitnami-labs/sealed-secrets)
- Ingress management with [nginx ingress controller](https://kubernetes.github.io/ingress-nginx/)
- Protective monitoring provided by GDS TechOps CyberSecurity with [Splunk](https://www.splunk.com/)
- Cloud infrastructure hosted on [AWS](https://aws.amazom.com) in three availability zones in the London region managed with [Terraform](https://www.terraform.io/)


## Help and support
For help or support:
- message the team on the GDS Supported Platform slack channel [#re-gsp](https://gds.slack.com/messages/CDA7YSP0D)
- raise an issue on[ gsp-team-manual](https://github.com/alphagov/gsp-team-manual/issues)
- email the GDS supported platform team re-GSP-team@digital.cabinet-office.gov.uk
- read our [documentation](/docs)
