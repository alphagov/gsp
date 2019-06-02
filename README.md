> This repository (gsp-terraform-ignition) will eventually become the main repository for the GDS Supported Platform. It will be renamed to something more sensible (e.g. gds-supported-platform, gsp, to de decided) however at the moment there are dependencies on this so it has to keep the same name.

# What is the GDS Supported Platform?

## Overview

The GDS Supported Platform is the strategic container platform for hosting services developed for the [Government Digital Service](https://www.gov.uk/government/organisations/government-digital-service).

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

## Running locally

It is possible to run a GSP cluster locally using [Minikube](https://kubernetes.io/docs/setup/minikube/).

To create a cluster run:

```
./scripts/gsp-local.sh create
```

After more than 5 minutes, but less than 10 the script should terminate and a cluster will be available:

```
kubectl cluster-info
```

You can then connect to some of the GSP-provided services:

```
# Workaround the fact Istio doesn't let you have specify ports on
# hosts with VirtualServices so we have to run it on a priviledged
# port:
https://github.com/istio/istio/issues/6469
sudo --preserve-env kubectl port-forward service/istio-ingressgateway -n istio-system 80:80
```
- Open a browser
- Navigate to `http://registry.local.govsandbox.uk`

When you're finished you should destroy it:

```
./scripts/gsp-local.sh destroy
```

## Help and support
For help or support:
- read our [documentation](/docs)
- raise an [issue](https://github.com/alphagov/gsp-terraform-ignition/issues)
- message the team on the Reliability [Engineering Slack channel](https://gds.slack.com/messages/CAD6NP598) [#reliability-eng](https://gds.slack.com/messages/CAD6NP598)
