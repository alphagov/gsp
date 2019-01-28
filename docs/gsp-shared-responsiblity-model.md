# GDS Supported Platform (GSP) shared responsibility

TechOps (including Reliability Engineering, Cyber Security and User Support) uses a [shared responsibility model](https://aws.amazon.com/compliance/shared-responsibility-model/) to provide a supported platform to GDS product teams. Broadly, TechOps builds, runs and maintains the infrastructure, while the product teams maintain the applications that they wish to run on the platform.

TechOps is only responsible for internal GDS product teams who host their service on the [GSP](https://github.com/alphagov/gsp-team-manual).

## TechOps' responsibilities

If you experience problems with the GSP, contact us using the [#techops Slack channel](https://gds.slack.com/messages/CBBK73196).

[something probably needs to go in here about Zendesk?]

### Provisioning Amazon Web Services (AWS) accounts and infrastructure

When a team has decided to use the GSP for hosting their application or service, Reliability Engineering will create an AWS account for the team (unless one exists for their umbrella programme) and provision the underlying infrastructure for teams to deploy their app.

### Updates and upgrades

When infrastructural components need to be upgraded/updated, Reliability Engineering (RE) will do so for each team's infrastructure. This may happen often and will rarely impact the service or their users, so RE will not notify the service team before performing the upgrade/update unless there will be a discernible impact.

### Logging, monitoring and alerting

TechOps will provide logging, monitoring and alerting for the GSP. We will work with teams to ensure that it is setup in the most useful way for teams. Specifically:

RE will ship logs to CloudWatch for further distribution. These logs will include:

- [insert log types here]

Cyber Security will then work with teams to establish use cases for the [Cloud Security Watch](#link) service that they provide. Cloud Security Watch will alert service teams when these use cases are triggered. This alert will include helpful information as to how the service team should deal with the situation.

If the service team is unable to respond to the situation themselves, they should contact the Cyber Security team either [on Slack](https://gds.slack.com/messages/CCMPJKFDK) or through the [out-of-hours procedure](#link), if necessary.

Cyber Security will also provide access to [Splunk](https://www.splunk.com/) and training for its use by service teams. Service teams will then be able to use Splunk for further monitoring.

If, when dealing with service issues, security events/incidents, etc., service teams identify infrastructural issues, TechOps (RE or Cyber Security, depending on where the issue occurs) will work with the service team to rectify the issue.

Due to the way in which applications are deployed and the production of logs occurs, TechOps will not be responsible for the structuring of logs or consistent identification of log streams (see *Logging* section in *Product team responsibilities* below).

TechOps will provide monitoring, using [Prometheus](), and alerting, using [AlertManager]() and [PagerDuty]().

TechOps will work with teams to identify the most valuable things to monitor and alert on (including the above Cloud Security Watch offering) and then to implement the correct logging, monitoring and alerting to facilitate this.

TechOps will setup monitoring and alerting to ensure that we provide reliable and secure infrastructure, in line with our [Service Level Objectives]().

[I think maybe this needs to be split out into separate "logging" and "monitoring and alerting" sections, but I don't know how...]

### Cyber Security and acting on vulnerabilities

Reliability Engineering will monitor [CVEs](https://en.wikipedia.org/wiki/Common_Vulnerabilities_and_Exposures) that relate to the technologies that underpin the infrastructure that we provide. We will then mitigate any identified vulnerabilities with [THIS PERIOD OF TIME].

Cyber Security will provide tooling and advice to teams, in order that they can identify and act upon security events and incidents. Tooling includes [Cloud Security Watch](https://github.com/alphagov/csw-infra) and [Splunk](https://www.splunk.com/). Advice includes working with teams to identify possible vulnerabilities that need to be monitored for exploitation, and then helping to set up that monitoring/alerting.

### User management

TechOps will control access to the underlying infrastructure using [AWS]() [IAM roles](). [some other stuff here]

## Product team responsibilities

### Updates and upgrades

Product teams will be responsible for any updates/upgrades to the application and any upstream dependencies (such as libraries and packages).

TechOps will, however, make this as easy as possible to implement, as changes such as these can be made live by making the changes in the code base before merging them to the master branch of the repository.

### Logging

Product teams are responsible for ensuring that they configure their applications and Docker images to ship logs [link to guidance?].

They must also ensure that their logs have a sensible and logical structure, to make them easy to identify and filter. For example:

[insert example of logically structured logs or link to guidance]

Similarly, product teams must name their containers sensibly and logically for consistent identification of log streams. For example:

[insert example of logically named containers or link to guidance]

### Monitoring and responding to alerts

It is the product team's responsibility to monitor their product/service (using tooling provided by TechOps - see above) and to respond to both monitoring and alerts as required.

The product team is also expected to allocate time to defining Service Level Indicators (SLIs) and Objectives (SLOs), in order that we can collectively setup the correct monitoring and alerting.

### User management

You must ensure that you control access to GitHub. As any code that is merged to the master branch of a repository will automatically be deployed, it is vital that only [trusted individuals](link to GDS Way bit about SC clearance and stuff) are able to merge pull requests.

### [What of the following (if any) will we do vs. what the service team should do?]

You must secure your AWS infrastructure. For example:

* control egress traffic by [implementing VPC egress controls](https://aws.amazon.com/answers/networking/controlling-vpc-egress-traffic/)
* use [TLS and other secure protocols](https://www.ncsc.gov.uk/guidance/tls-external-facing-services) to protect data
* use secure coding practices like peer review
* secure developer machines - ask the GDS IT team for guidance
* manage secrets to [secure the build and deployment pipeline](https://www.ncsc.gov.uk/guidance/secure-build-and-deployment-pipeline)
* control access to provisioned machines and other AWS services
* implement protective monitoring and [logging for security purposes](https://www.ncsc.gov.uk/guidance/introduction-logging-security-purposes)
* use [security hardening for VMs](https://gds-way.cloudapps.digital/standards/operating-systems.html)
* enabling infrastructure monitoring using [AWS CloudWatch](https://aws.amazon.com/cloudwatch/)

You must [transfer ownership](https://github.com/alphagov/re-build-systems/blob/master/examples/gds_specific_dns_and_jenkins/README.md#provision-the-main-jenkins-infrastructure) of the OAuth app to the [GitHub alphagov organisation](https://github.com/alphagov) once you have provisioned Jenkins. This prevents unauthorised access to the build system if the owner of the OAuth app leaves GDS.

## Third party responsibilities

Third party providers are responsible for making sure their systems are updated and available.

### Amazon

Amazon maintains the AWS infrastructure and is responsible for updating it. Reliability Engineering and product teams do not need to update or upgrade AWS.

### Docker.io

When you provision or reprovision your infrastructure you use the latest version of Docker. Docker is responsible for the maintenance and support of new and current versions.

### [any others explicitly?]
