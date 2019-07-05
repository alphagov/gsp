# ADR001: Support Model

## Status

Accepted

## Context

How we structure our support model has a large impact on the design of any
systems we put in place or services we provide.

The different models could be considered on a spectrum:

* A **Platform as a Service** model:
    * A larger dedicated team of SREs provide a centralised full end-to-end
      service for building, deploying, running applications and required
      services for Service Teams to consume
    * Usually implies centralised infrastructure and benefits from the economies
      of scale that such an architecture brings
    * Requires little to no ops knowledge from Service Teams
    * Evolution of the system is dependent on the platform team
* A **shared responsibility** or **supported workflow** model:
    * A smaller team of dedicated SREs provide the lower-level infrastructure
      primitives, and – by regularly rotating/embedding/consulting within teams
      – work to identify common needs, and improve the re-usability and
      reliability of shareable solutions and maintain best-practice
      documentation
    * Makes no assumptions about what is "centralised"
    * Requires some knowledge of the ops side from service teams or a structure
      for support from the dedicated SRE team
    * Evolution of the system is a shared responsibility
* A **Paved Path** model:
    * Teams have their own dedicated SREs focused solely on their team's needs
      and knowledge sharing is achieved by regularly contributing to
      organisation best-practices, guidance documents or reusable code where and
      when possible
    * Usually implies de-centralised infrastructure with each team running their
      own
    * Requires strong knowledge of ops side within Service Team
    * Evolution of the system is dependent on the Service Team


The table below aims to illustrate where responsibility for the maintenance,
evolution and iteration of various parts of such a system lie in each of
these models.

| Responsibility  | PaaS Model | Supported Workflow Model | Paved Path Model |
|---|---|---|---|
|  | Give us your code, we'll keep it running for you | Use our tools and follow our guidance, and we'll support you | Here's some examples and guidance, do it yourself |
| Application lifecycle (HA, Capacity) | RE | RE/ServiceTeam | ServiceTeam |
| Observability etc | RE | RE/ServiceTeam | ServiceTeam |
| Backing services | RE | RE/ServiceTeam | ServiceTeam |
| Iterating new features | RE | RE/ServiceTeam | ServiceTeam |
| Workflow/Process | RE | RE | ServiceTeam |
| Productionising common features | RE | RE | ServiceTeam |
| Infrastructure Responsibility | RE | RE | ServiceTeam |

Traditionally GDS teams were operating closest to a "paved path" model where teams took complete control over their own infrastructure. GDS also operates the [GOV.UK PaaS](https://cloud.service.gov.uk), a platform targeting other government departments with limited operations capability.

More recently GDS has created a Reliability Engineering team (RE) to provide support for common needs of Service Teams and reduce duplication of effort around areas like infrastructure, build tooling and observability.

With more of the operations skills centralised in Reliability Engineering, models on the "paved path" end of the spectrum become more difficult to sustain. This leaves us with two choices:

* provide a full-stack PaaS model, growing the size of the PaaS team as required to keep up with demand for new features.
* provide only the lower-level building blocks, guidance and let service teams contribute to the development of new features

## Decision

We will design a system around a shared responsibility between Reliability Engineering and Service Teams.

We believe that the middle ground where service teams have a high level of responsibility over their deployments but Reliability Engineering takes most of the responsibility for the lower-level deployment primitives as well as providing solutions to the common needs of service teams (such as deployment workflows, monitoring, deployment patterns) will maintain the flexibility service teams need when iterating new features as well as providing the framework for getting support from Reliability Engineering when such new features could benefit from hardening or sharing with the organisation as a whole.

* We should provide, promote and support the use of a common language for declaratively describing infrastructure and deployments
* We should provide, promote and support low-level deployment primitives suitable for service teams to build upon
* We should work closely and continuously with service teams to reduce duplication of effort by identify common needs and providing reliable, reusable solutions

## Consequences

* Implies a framework for embedding/loaning/floating SREs within Service Teams so that knowledge can be shared and common needs can be identified and improved upon
* Will require a migration to a common declarative deployment language
* Gives more control/flexibility to Service Teams than a PaaS model at the cost of potentially more complex deployment patterns
* Keeping Service Teams converging on a smaller set of best-practice solutions would be harder than with a centrally enforced PaaS model
