# Architecture Decision Records

We document our decisions using [Architecture Decision Records](https://github.com/alphagov/gsp-team-manual/tree/master/adr) as recommended by the [GDS Way](https://gds-way.cloudapps.digital/standards/architecture-decisions.html)

## Index

- [ADR #000 - Decision Record Template](ADR000-template.md)
- [ADR #001 - Support Model](ADR001-support-model.md)
- [ADR #002 - Containers](ADR002-containers.md)
- [ADR #003 - Container Orchestration](ADR003-container-orchestration.md)
- [ADR #004 - Tenant Isolation](ADR004-tenant-isolation.md) [superseded by [ADR #024 - Soft Multi-tenancy](ADR024-soft-multitenancy.md)]
- [ADR #005 - Ingress](ADR005-ingress.md) [superseded by [ADR #025 - Ingress](ADR025-ingress.md)]
- [ADR #006 - Cluster Authentication Method](ADR006-authentication-method.md) [superseded by [ADR #023 - Cluster Authentication](ADR023-cluster-authentication.md)]
- [ADR #007 - Identity Provider](ADR007-identity-provider.md) [superseded by [ADR #023 - Cluster Authentication](ADR023-cluster-authentication.md)]
- [ADR #008 - Continuous delivery workflow](ADR008-continuous-delivery-workflow.md)
- [ADR #009 - Multi-tenancy for CI and CD](ADR009-multitenant-ci-cd.md) [superseded by [ADR #029 - Pull based Continuous Delivery tools](ADR029-continuous-delivery-tools.md)]
- [ADR #010 - Placement of CI and CD Tools](ADR010-placement-of-ci-cd-tools.md) [superseded by [ADR #029 - Pull based Continuous Delivery tools](ADR029-continuous-delivery-tools.md)]
- [ADR #011 - Build Artefacts](ADR011-build-artefacts.md)
- [ADR #012 - Docker image repositories](ADR012-docker-image-repositories.md) [superseded by [ADR #028 - Container Tools](ADR028-container-tools.md)]
- [ADR #013 - CI & CD Tools](ADR013-ci-cd-tools.md)
- [ADR #014 - Sealed Secrets](ADR014-sealed-secrets.md)
- [ADR #015 - AWS IAM Authentication (for admins)](ADR015-aws-iam-authentication.md) [superseded by [ADR #023 - Cluster Authentication](ADR023-cluster-authentication.md)]
- [ADR #016 - Code verification](ADR016-code-verification.md)
- [ADR #017 - Vendor provided container orchestration](ADR017-vendor-provided-container-orchestration.md)
- [ADR #018 - Local Development Environment](ADR018-local-development.md)
- [ADR #019 - Service Mesh](ADR019-service-mesh.md)
- [ADR #020 - Metrics](ADR020-metrics.md)
- [ADR #021 - Alerting](ADR021-alerting.md)
- [ADR #022 - Logging](ADR022-logging.md)
- [ADR #023 - Cluster Authentication](ADR023-cluster-authentication.md)
- [ADR #024 - Soft Multi-tenancy](ADR024-soft-multitenancy.md)
- [ADR #025 - Ingress](ADR025-ingress.md)
- [ADR #028 - Container Tools](ADR028-container-tools.md)
- [ADR #029 - Pull based Continuous delivery tools](ADR029-continuous-delivery-tools.md)
- [ADR #030 - AWS Service Operator](ADR030-aws-service-operator.md)
- [ADR #031 - Postgres](ADR031-postgres.md)
- [ADR #032 - SRE permissions](ADR032-sre-permissions.md)
- [ADR #033 - NLB for mTLS](ADR033-nlb-for-mtls.md)
- [ADR #034 - Single service operator](ADR034-one-service-operator-different-resource-kinds.md)
- [ADR #035 - Aurora postgres](ADR035-aurora-postgres.md)
- [ADR #036 - CloudHSM isolation](ADR036-hsm-isolation-in-detail.md)
- [ADR #037 - Per namespace istio gateways](ADR037-per-namespace-gateways.md)
- [ADR #038 - SRE Permissions for Istio](ADR038-sre-permissions-istio.md)
- [ADR #039 - Restricting CloudHSM network access to particular namespaces](ADR039-cloudhsm-namespace-network-policy.md)
- [ADR #040 - Ensuring cluster stability while replacing nodes](ADR040-cluster-stability-node-replacement.md)
- [ADR #041 - Service operator provisioned policies](ADR041-service-operated-policies.md)
- [ADR #042 - Static ingress IP workaround](ADR042-static-ingress-ip-workaround.md)
- [ADR #043 - Kubernetes resource access rules](ADR043-k8s-resource-access.md)
- [ADR #044 - Security improvements](ADR044-security-improvements.md)
- [ADR #045 - Dev namespaces](ADR045-dev-namespaces.md)
