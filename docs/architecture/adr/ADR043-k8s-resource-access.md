# ADR043: Kubernetes resource access rules

## Status

Accepted

## Context

Several different levels of access are required within a kubernetes cluster. The
GSP uses role based access control so these levels are granted to users and
groups via roles.

## Decision

We will create two levels of access within each namespace:

* Operator
* Auditor

The Operator is a relatively permissive read-write role within the namespace.
Developers working on branches that are not part of the release process may be
granted this role in certain namespaces. This is also the role the in-cluster
concourse team for each namespace will be granted.

The Auditor will be given to all authenticated users in the cluster. This should
allow for debugging of issues and incidents and basic remedial actions without
needing formal escalation procedures.

The complete list of resource permissions is given in Appendix A.

## Consequences

* Site Reliability Engineers (SRE) will be able to effectively operate the GSP
  while reducing the need to escalate to "admin" level access.
* service team members (e.g. devs) will be able to iterate changes quickly via
  branch-based development and/or `kubectl apply` level access to specific
  namespaces.
* service team members will not be able to interfere with other devs'
  namespaces, "system" namespaces or "prod" namespaces.
* SREs and service team members will not be able to cause major production
  service outages. Degradation will be possible through deletion of certain
  resources but that will be easily fixable by running the appropriate
  deployment pipeline.


## Appendix A: Permissions rules

---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: auditor
rules:
- apiGroups: [""]
  resources:
  - configmaps
  - endpoints
  - limitranges
  - pods
  - podtemplates
  - resourcequotas
  - serviceaccounts
  - services
  verbs:
  - delete
  - get
  - list
  - watch
- apiGroups: [""]
  resources:
  - events
  - namespaces
  - nodes
  - persistentvolumeclaims
  - persistentvolumes
  - pods/log
  verbs:
  - get
  - list
  - watch
- apiGroups: [""]
  resources:
  - secrets
  verbs:
  - delete
  - list
- apiGroups: [""]
  resources:
  - pods/portforward
  verbs:
  - create
- apiGroups: [""]
  resources:
  - services/proxy
  verbs:
  - get

- apiGroups: ["access.govsvc.uk"]
  resources:
  - principals
  verbs:
  - delete
  - get
  - list
  - watch

- apiGroups: ["admissionregistration.k8s.io"]
  resources:
  - mutatingwebhookconfigurations
  - validatingwebhookconfigurations
  verbs:
  - get
  - list
  - watch

- apiGroups: ["apiextensions.k8s.io"]
  resources:
  - customresourcedefinitions
  verbs:
  - get
  - list
  - watch

- apiGroups: ["apiregistration.k8s.io"]
  resources:
  - apiservices
  verbs:
  - get
  - list
  - watch

- apiGroups: ["apps"]
  resources:
  - daemonsets
  verbs:
  - get
  - list
  - watch
- apiGroups: ["apps"]
  resources:
  - deployments
  - replicasets
  - statefulsets
  verbs:
  - delete
  - get
  - list
  - watch

- apiGroups: ["authentication.istio.io"]
  resources:
  - meshpolicies
  - policies
  verbs:
  - get
  - list
  - watch

- apiGroups: ["autoscaling"]
  resources:
  - horizontalpodautoscalers
  verbs:
  - delete
  - get
  - list
  - watch

- apiGroups: ["batch"]
  resources:
  - cronjobs
  - jobs
  verbs:
  - delete
  - get
  - list
  - watch

- apiGroups: ["bitnami.com"]
  resources:
  - sealedsecrets
  verbs:
  - delete
  - get
  - list
  - watch

- apiGroups: ["cert-manager.io"]
  resources:
  - certificaterequests
  - certificates
  - issuers
  verbs:
  - delete
  - get
  - list
  - watch
- apiGroups: ["cert-manager.io"]
  resources:
  - clusterissuers
  verbs:
  - get
  - list
  - watch

- apiGroups: ["certificates.k8s.io"]
  resources:
  - certificatesigningrequests
  verbs:
  - delete
  - get
  - list
  - watch

- apiGroups: ["concourse.k8s.io"]
  resources:
  - pipelines
  verbs:
  - delete
  - get
  - list
  - watch
- apiGroups: ["concourse.k8s.io"]
  resources:
  - teams
  verbs:
  - get
  - list
  - watch

- apiGroups: ["config.gatekeeper.sh"]
  resources:
  - configs
  verbs:
  - get
  - list
  - watch

- apiGroups: ["config.istio.io"]
  resources:
  - adapters
  - apikeys
  - attributemanifests
  - authorizations
  - circonuses
  - deniers
  - dogstatsds
  - edges
  - fluentds
  - httpapispecbindings
  - httpapispecs
  - kuberneteses
  - listcheckers
  - listentries
  - logentries
  - noops
  - opas
  - quotaspecbindings
  - quotaspecs
  - rbacs
  - reportnothings
  - solarwindses
  - stackdrivers
  - statsds
  - zipkins
  verbs:
  - get
  - list
  - watch
- apiGroups: ["config.istio.io"]
  resources:
  - bypasses
  - checknothings
  - cloudwatches
  - instances
  - handlers
  - kubernetesenvs
  - memquotas
  - metrics
  - prometheuses
  - quotas
  - redisquotas
  - rules
  - signalfxs
  - stdios
  - templates
  - tracespans
  verbs:
  - delete
  - get
  - list
  - watch

- apiGroups: ["coordination.k8s.io"]
  resources:
  - leases
  verbs:
  - delete
  - get
  - list
  - watch

- apiGroups: ["crd.k8s.amazonaws.com"]
  resources:
  - eniconfigs
  verbs:
  - get
  - list
  - watch

- apiGroups: ["crd.projectcalico.org"]
  resources:
  - bgpconfigurations
  - bgppeers
  - blockaffinities
  - clusterinformations
  - felixconfigurations
  - globalnetworkpolicies
  - globalnetworksets
  - hostendpoints
  - ipamblocks
  - ippools
  - networksets
  verbs:
  - get
  - list
  - watch
- apiGroups: ["crd.projectcalico.org"]
  resources:
  - networkpolicies
  verbs:
  - delete
  - get
  - list
  - watch

- apiGroups: ["database.govsvc.uk"]
  resources:
  - postgres
  verbs:
  - delete
  - get
  - list
  - watch

- apiGroups: ["events.k8s.io"]
  resources:
  - events
  verbs:
  - get
  - list
  - watch

- apiGroups: ["extensions"]
  resources:
  - daemonsets
  - ingresses
  - podsecuritypolicies
  verbs:
  - get
  - list
  - watch
- apiGroups: ["extensions"]
  resources:
  - deployments
  - networkpolicies
  - replicasets
  verbs:
  - delete
  - get
  - list
  - watch

- apiGroups: ["metrics.k8s.io"]
  resources:
  - nodes
  - pods
  verbs:
  - get
  - list
  - watch

- apiGroups: ["monitoring.coreos.com"]
  resources:
  - alertmanagers
  - prometheuses
  verbs:
  - get
  - list
  - watch
- apiGroups: ["monitoring.coreos.com"]
  resources:
  - podmonitors
  - prometheusrules
  - servicemonitors
  verbs:
  - delete
  - get
  - list
  - watch

- apiGroups: ["networking.istio.io"]
  resources:
  - destinationrules
  - envoyfilters
  - gateways
  - serviceentries
  - sidecars
  - virtualservices
  verbs:
  - delete
  - get
  - list
  - watch

- apiGroups: ["networking.k8s.io"]
  resources:
  - ingresses
  verbs:
  - get
  - list
  - watch
- apiGroups: ["networking.k8s.io"]
  resources:
  - networkpolicies
  verbs:
  - delete
  - get
  - list
  - watch

- apiGroups: ["node.k8s.io"]
  resources:
  - runtimeclasses
  verbs:
  - get
  - list
  - watch

- apiGroups: ["policy"]
  resources:
  - poddisruptionbudgets
  verbs:
  - delete
  - get
  - list
  - watch
- apiGroups: ["policy"]
  resources:
  - podsecuritypolicies
  verbs:
  - get
  - list
  - watch

- apiGroups: ["queue.govsvc.uk"]
  resources:
  - sqs
  verbs:
  - delete
  - get
  - list
  - watch

- apiGroups: ["rbac.authorization.k8s.io"]
  resources:
  - clusterrolebindings
  - clusterroles
  verbs:
  - get
  - list
  - watch
- apiGroups: ["rbac.authorization.k8s.io"]
  resources:
  - rolebindings
  - roles
  verbs:
  - delete
  - get
  - list
  - watch

- apiGroups: ["rbac.istio.io"]
  resources:
  - authorizationpolicies
  - clusterrbacconfigs
  - rbacconfigs
  - servicerolebindings
  - serviceroles
  verbs:
  - get
  - list
  - watch

- apiGroups: ["scheduling.k8s.io"]
  resources:
  - priorityclasses
  verbs:
  - get
  - list
  - watch

- apiGroups: ["storage.govsvc.uk"]
  resources:
  - s3buckets
  verbs:
  - delete
  - get
  - list
  - watch

- apiGroups: ["storage.k8s.io"]
  resources:
  - csidrivers
  - csinodes
  - storageclasses
  - volumeattachments
  verbs:
  - get
  - list
  - watch

- apiGroups: ["templates.gatekeeper.sh"]
  resources:
  - constrainttemplates
  verbs:
  - get
  - list
  - watch

- apiGroups: ["verify.gov.uk"]
  resources:
  - certificaterequests
  - metadata
  verbs:
  - delete
  - get
  - list
  - watch

- apiGroups: ["webhook.cert-manager.io"]
  resources:
  - mutations
  - validations
  verbs:
  - get
  - list
  - watch

---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: operator
rules:
- apiGroups: [""]
  resources:
  - configmaps
  - endpoints
  - limitranges
  - pods
  - podtemplates
  - resourcequotas
  - serviceaccounts
  - services
  verbs:
  - create
  - patch
  - update
- apiGroups: [""]
  resources:
  - persistentvolumeclaims
  - persistentvolumes
  verbs:
  - create
  - delete
  - patch
  - update
- apiGroups: [""]
  resources:
  - secrets
  verbs:
  - create
  - get
  - patch
  - update
  - watch
- apiGroups: [""]
  resources:
  - pods/exec
  verbs:
  - create

- apiGroups: ["access.govsvc.uk"]
  resources:
  - principals
  verbs:
  - create
  - patch
  - update

- apiGroups: ["apps"]
  resources:
  - deployments
  - replicasets
  - statefulsets
  verbs:
  - create
  - patch
  - update

- apiGroups: ["authentication.istio.io"]
  resources:
  - policies
  verbs:
  - create
  - delete
  - patch
  - update

- apiGroups: ["autoscaling"]
  resources:
  - horizontalpodautoscalers
  verbs:
  - create
  - patch
  - update

- apiGroups: ["batch"]
  resources:
  - cronjobs
  - jobs
  verbs:
  - create
  - patch
  - update

- apiGroups: ["bitnami.com"]
  resources:
  - sealedsecrets
  verbs:
  - create
  - patch
  - update

- apiGroups: ["cert-manager.io"]
  resources:
  - certificaterequests
  - certificates
  - issuers
  verbs:
  - create
  - patch
  - update

- apiGroups: ["certificates.k8s.io"]
  resources:
  - certificatesigningrequests
  - pipelines
  verbs:
  - create
  - patch
  - update

- apiGroups: ["config.istio.io"]
  resources:
  - bypasses
  - checknothings
  - cloudwatches
  - handlers
  - instances
  - kubernetesenvs
  - memquotas
  - metrics
  - prometheuses
  - quotas
  - redisquotas
  - rules
  - signalfxs
  - stdios
  - templates
  - tracespans
  verbs:
  - create
  - patch
  - update

- apiGroups: ["coordination.k8s.io"]
  resources:
  - leases
  verbs:
  - create
  - patch
  - update

- apiGroups: ["crd.projectcalico.org"]
  resources:
  - networkpolicies
  verbs:
  - create
  - patch
  - update

- apiGroups: ["database.govsvc.uk"]
  resources:
  - postgres
  verbs:
  - create
  - patch
  - update

- apiGroups: ["extensions"]
  resources:
  - deployments
  - networkpolicies
  - replicasets
  verbs:
  - create
  - patch
  - update

- apiGroups: ["monitoring.coreos.com"]
  resources:
  - podmonitors
  - prometheusrules
  - servicemonitors
  verbs:
  - create
  - patch
  - update

- apiGroups: ["networking.istio.io"]
  resources:
  - destinationrules
  - envoyfilters
  - gateways
  - serviceentries
  - sidecars
  - virtualservices
  verbs:
  - create
  - patch
  - update

- apiGroups: ["networking.k8s.io"]
  resources:
  - networkpolicies
  verbs:
  - create
  - patch
  - update

- apiGroups: ["policy"]
  resources:
  - poddisruptionbudgets
  verbs:
  - create
  - patch
  - update

- apiGroups: ["queue.govsvc.uk"]
  resources:
  - sqs
  verbs:
  - create
  - patch
  - update

- apiGroups: ["rbac.authorization.k8s.io"]
  resources:
  - rolebindings
  - roles
  verbs:
  - create
  - patch
  - update

- apiGroups: ["storage.govsvc.uk"]
  resources:
  - s3buckets
  verbs:
  - create
  - patch
  - update

- apiGroups: ["verify.gov.uk"]
  resources:
  - certificaterequests
  - metadata
  verbs:
  - create
  - patch
  - update
