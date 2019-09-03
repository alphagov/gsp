# Sealing secrets

Create a regular Kubernetes `Secret`:

```
# lol.yaml
apiVersion: v1
kind: Secret
metadata:
  name: lol
  namespace: your-namespace
type: Opaque
data:
  # base64 encode the values
  lol: bG9sCg==
```

Convert the `Secret` to a `SealedSecret`:

```
gds sandbox seal < lol.yaml > lol-sealed-secret.yaml --format yaml
```

Commit the `SealedSecret`.
