Making images private in Harbor
===============================

Sometimes, we need to make our Docker images private. By default, the GSP Harbor
is publicly accessible. Here are instructions on how to make images
private.

1. Add the following to your build-pipeline.yml:

```
  harbor_source: &harbor_source
    username: ((harbor.harbor_username))
    password: ((harbor.harbor_password))
    harbor:
      url: ((harbor.harbor_url))
      prevent_vul: "false"
      public: "false"
    notary:
      url: ((harbor.notary_url))
      root_key: ((harbor.root_key))
      delegate_key: ((harbor.ci_key))
      passphrase:
        root: ((harbor.notary_root_passphrase))
        snapshot: ((harbor.notary_snapshot_passphrase))
        targets: ((harbor.notary_targets_passphrase))
        delegation: ((harbor.notary_delegation_passphrase))

  resource_types:
  - name: harbor
    type: docker-image
    privileged: true
    source:
      repository: govsvc/gsp-harbor-docker-image-resource
      tag: "0.0.1553882420"

  resources:
  - name: my-image
    type: harbor
    source:
      <<: *harbor_source
      repository: registry.((cluster.domain))/my-collection/my-image
```

The crucial line is `public: "false"` to ensure that anonymous users (ie
public users) can't pull images from that Harbor repo.

2. Add `imagePullSecrets` to your Helm chart in order to use the
   builtin secret `registry-creds` for pulling the private Docker
   image:

```
template:
    metadata:
      labels:
        app.kubernetes.io/name: gateway
        app.kubernetes.io/instance: {{ .Release.Name }}
    spec:
      restartPolicy: Always
      volumes:
        ...
      containers:
        ...
      imagePullSecrets:
      - name: registry-creds
```
