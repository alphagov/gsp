# Permissions

The permissions a given user will have in a given namespace in all clusters in a given account
depend on their configuration in [gds-trusted-developers]. Specifically, the
`roles` they have.

## Cluster admins

A `gds-trusted-developer` may be configured as a cluster admin in all clusters
in an account. The cluster admin permissions should be used only during periods
where it is strictly necessary (such as during an incident) and should be
surrendered immediately following the return to normal service.

To configure a user as a cluster admin:

```
roles:
- account: verify
  role: admin
```

## Namespace operators

An "operator" in a namespace has an elevated set of permissions to accelerate
the feedback cycle of development for a tenant. For example they are able to
create arbitrary, namespace-scoped resources using `kubectl`, read secrets, and
view or edit pipelines. To elevate a `gds-trusted-developer` to an "operator" in
a given namespace:

```
roles:
- account: verify
  role: operator
  namespace: verify-my-dev-namespace
```

## Cluster auditors

All `gds-trusted-developers` in a given account are given "auditor" access to
all clusters in the account. This gives basic read access to the whole cluster (except for some
sensitive resources such as secrets).

To configure an "auditor":

```
roles:
- account: verify
  role: auditor
```

## Further info

* [ADR043] k8s resource access
* [ADR044] security improvements
* [ADR045] dev namespaces


[ADR043]:
https://github.com/alphagov/gsp/blob/master/docs/architecture/adr/ADR043-k8s-resource-access.md
[ADR044]:
https://github.com/alphagov/gsp/blob/master/docs/architecture/adr/ADR044-security-improvements.md
[ADR045]:
https://github.com/alphagov/gsp/blob/master/docs/architecture/adr/ADR045-dev-namespaces.md
[gds-trusted-developers]: https://github.com/alphagov/gds-trusted-developers
