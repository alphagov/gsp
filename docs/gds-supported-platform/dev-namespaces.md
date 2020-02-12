# Dev namespaces

A dev namespace is the same as any other tenant namespace: it is created by the
GSP as a result of the configuration in the `*-cluster-config` repository.

For example:

```
namespaces:
- name: verify-my-dev-namespace
  owner: alphagov
  repository: doc-checking
  branch: my-dev-branch
  path: ci/dev
  requiredApprovalCount: 0
  ingress:
    enabled: true
```

To gain the elevated permissions in this namespace a user from
`gds-trusted-developers` needs to be granted the `operator` role.

```
name: jeff.jeffersonreference
email: jeff.jefferson@email.com
ARN: arn:aws:iam::000000000006:user/jeff.jefferson@email.com
roles:
- account: verify
  role: operator
  namespace: verify-my-dev-namespace
hardware:
  id: 9000005
  type: yubikey
github: jefferson678
```

For details on what additional permissions `jeff.jeffersonreference` would now
have in that namespace see [the operator role
definition](https://github.com/alphagov/gsp/blob/master/charts/gsp-cluster/templates/00-aws-auth/operator-cluster-role.yaml).
