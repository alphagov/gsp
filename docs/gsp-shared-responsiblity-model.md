# GDS Supported Platform (GSP) shared responsibility

TechOps (including Reliability Engineering, Cyber Security and User Support) uses a [shared responsibility model](https://aws.amazon.com/compliance/shared-responsibility-model/) to provide a supported platform to GDS product teams. Broadly, TechOps builds, runs and maintains the infrastructure and supporting services, while the product teams maintain the applications that they wish to run on the platform.

TechOps is only responsible for internal GDS product teams who host their service on the [GDS Supported Platform](https://github.com/alphagov/gsp-team-manual) and the [GOV.UK Platform as a Service (PaaS)](https://www.cloud.service.gov.uk/) and its tenants.

## Summary of responsibilities

| Activity | Product team | Reliability Engineering | Cyber Security |
| -------- | -------- | -------- | -------- |
| Provide account (e.g. AWS, PagerDuty, Alert Manager)     |  | ✓ |  |
| Provide account (e.g. Splunk)     |  |  | ✓ |  
| Create build/deploy pipeline     | ✓ | ✓* |  |
| Respond to incidents     | ✓ | ✓** | ✓** |
| Respond to other alerts     | ✓ | ✓** |  |
| Set up infrastructure |  | ✓ |  |
| Maintain infrastructure     |  | ✓ |  |
| Maintain product/application(s) | ✓ |  |  |
| Set up & structure logging | ✓ | ✓* | ✓* |
| Provide tooling for logging |  | ✓ | ✓ |
| Monitor the reliability of key user journeys (SLIs) | ✓ | ✓* |  |
| Set shared reliability goals for each SLI (SLOs) |✓ | ✓ |  |
| Agree on policy for breaking SLOs |✓ | ✓ |  |

\* = in a supporting role/through documentation
** = if called upon / if the main team can’t deal with the situation

## TechOps' responsibilities

If you (a product team) experience problems with the GDS Supported Platform, contact us using the [#techops](https://gds.slack.com/messages/CBBK73196) Slack channel.

Contact Cyber Security using the [#cyber-security-help](https://gds.slack.com/messages/CCMPJKFDK) Slack channel.

### Provisioning Amazon Web Services (AWS) accounts and infrastructure

When a team has decided to use the GDS Supported Platform for hosting their application or service, Reliability Engineering (RE) will create an [AWS account for the team](https://reliability-engineering.cloudapps.digital/iaas.html#create-aws-accounts) (unless one exists for their wider programme, for example, GOV.UK or GOV.UK Verify) and provision the underlying infrastructure for teams to deploy their app.

### Updates and upgrades

When infrastructural components need to be upgraded/updated, RE will do so for each team's infrastructure. This may happen often and will rarely impact the service or their users, so RE will not notify the service team before performing the upgrade/update unless there will be a discernible impact.

### Setting up a build/deploy pipeline

Reliability Engineering will help teams to set up a build and deployment pipeline (otherwise known as a continuous integration and deployment (CI/CD) pipeline) as part of helping teams migrate to either the PaaS or the GSP.

This includes providing the tooling and documentation for product teams, as well as consulting/working with teams to get an initial pipeline setup.

RE will always be available on a consultant basis, should product teams need help.

RE is also responsible for providing tooling to monitor the pipeline, for example to monitor the progress (and success/failure) of jobs that are running.

### Logging, monitoring and alerting

TechOps will provide logging, monitoring and alerting for the GDS Supported Platform. We will work with teams to ensure that it is set up in the most useful way for teams.

RE will ship logs to CloudWatch for further distribution (for example, to Splunk).

Cyber Security will then work with teams to establish use cases for teams' protective monitoring using Splunk. Splunk will alert service teams when those use cases are triggered. Playbooks will include helpful information as to how the service team should deal with a given situation.

If the service team is unable to respond to the situation themselves or if it is categorised as a security incident, they should contact the Cyber Security team either on Slack or through the out-of-hours procedure, if necessary.

Cyber Security will also provide access to Splunk. Service teams will then be able to use Splunk to maintain their current protective monitoring use cases and develop new ones over time.

If, when dealing with service issues, security events/incidents, etc., service teams identify infrastructural issues, TechOps (RE or Cyber Security, depending on where the issue occurs) will work with the service team to resolve the issue.

Due to the way in which applications are deployed and the production of logs occurs, TechOps will not be responsible for the structuring of logs or consistent identification of log streams (see Logging section in Product team responsibilities below).

TechOps will provide:

- health and reliability monitoring using [Prometheus](https://prometheus.io/)
- alerting using [Prometheus' AlertManager](https://prometheus.io/docs/alerting/alertmanager/) and [PagerDuty](https://www.pagerduty.com/).

TechOps will work with teams to identify the most valuable things to monitor and alert on (including the protective monitoring offering) and then work together to implement the correct logging, monitoring and alerting to facilitate this.

TechOps will set up health and reliability monitoring and alerting to ensure that we provide reliable and secure infrastructure, in line with our service level objectives (SLOs).

### Ensuring the integrity of deployments

TechOps will provide tooling and enforcement of measures that ensure that the right code gets deployed to production in the right way.

Specifically, TechOps will enforce the established best practice of "two-eyes" (having at least two people look at the code before it is merged to the master branch of a repository), which will mean that no deployments can be made without being reviewed and signed by at least two authorised people.

### Responding to security events and incidents and acting on vulnerabilities

Reliability Engineering (RE) will monitor CVEs that relate to the technologies that underpin the infrastructure that we provide. RE will then work to mitigate any identified vulnerabilities that affect the platform as quickly as possible.

Cyber Security will provide tooling and advice to teams, in order that they can identify and act upon security events and incidents.

Tooling includes:

- [Cloud Security Watch](https://github.com/alphagov/csw-backend), a tool to detect misconfigurations in AWS
- [Splunk](https://www.splunk.com/), a security information and event management [(SIEM)](https://en.wikipedia.org/wiki/Security_information_and_event_management) tool to enable protective monitoring

Advice includes:

- working with teams to conduct [threat modelling](https://www.owasp.org/index.php/Category:Threat_Modeling)
- working with teams to identify use cases to mitigate risks that need to be protectively monitored
- supporting teams to ingest data sources to Splunk
- helping to set up protective monitoring and alerting

### User management

TechOps will control access to the underlying infrastructure using AWS IAM roles. TechOps is therefore responsible for allowing and removing this access in a timely manner.

All changes to user access and to the platform go through a two-eyes process to ensure a robust and secure approach to user access control.

## Product team responsibilities

### Updates and upgrades

Product teams will be responsible for any updates/upgrades to the application and any upstream dependencies (such as libraries and packages).

TechOps will, however, make this as easy as possible to implement, as changes such as these can be put into production by making the changes in the code base before merging them to the master branch of the repository.

### Setting up and maintaining a build/deploy pipeline

Product teams will need to work with TechOps to set up the initial pipeline, including learning (through documentation that TechOps provides and through working together) how to do so themselves.

Once the initial pipeline is established, the product team will be responsible for making sure that tests and specific procedures (for example. linting, formatting, promoting) are maintained as the application(s) develop(s).

Product teams are responsible for monitoring their pipelines for successful or failing jobs and acting upon that information.

### Logging

Product teams are responsible for ensuring that they configure their applications and Docker images to ship logs.

TechOps will develop guidance on logging in the near future. Refer to the [Splunk guidelines](http://dev.splunk.com/view/logging/SP-CAAAFCK) for now.

Some brief examples can be seen on the [Splunk website](http://dev.splunk.com/view/logging/SP-CAAAFCM).

Similarly, product teams must name their containers sensibly and logically for consistent identification of log streams. Again, we will provide more detailed guidance soon.

### Monitoring and responding to alerts

It is the product team's responsibility to monitor their product/service (using tooling provided by TechOps) and to respond to both monitoring and alerts as required. If the product team is unable to respond to the situation themselves or if it is categorised as a security incident, they should contact the Cyber Security team either on Slack or through the out-of-hours procedure, if necessary.

Product teams are responsible for the development, maintenance and engineering of their use cases in Splunk. As their products develop, they should ensure their protective monitoring stays effective and relevant. Cyber Security will support teams where required.

The product team is also expected to allocate time to defining service level indicators (SLIs) and objectives (SLOs), so that we can collectively set up the correct monitoring and alerting.

### User management

Product teams must ensure that they control access to GitHub. As any code that is merged to the master branch of a repository will automatically be deployed, it is vital that only trusted individuals are able to merge pull requests.

## Third party responsibilities

Third party providers are responsible for making sure their systems are updated and available.

### Amazon

Amazon maintains the AWS infrastructure and is responsible for updating it. Reliability Engineering and product teams do not need to update or upgrade AWS.

### Docker.io



[What's the escalation path for when the inevitable arguments happen?]
[something probably needs to go in here about Zendesk?]
See Cyber Security’s page for more information about protective monitoring, our tools and services [link to more detailed guidance?]().
[insert log types]
[link to out of hours procedure]
