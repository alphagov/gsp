# Sealing secrets

In order to store secrets safely in a Git repository you can seal secrets using a tool called `kubeseal` which will generate an encrypted
kubernetes resource using a public key from your target cluster that can only be decrypted into it's target namespace by the target cluster.

## Requirements

You will need the following cli tools:

* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* [kubeseal](https://github.com/bitnami-labs/sealed-secrets/releases)
* [jq](https://stedolan.github.io/jq/download/)

We recommend you also have:

* [aws-vault](https://github.com/99designs/aws-vault)
* a `KUBECONFIG=<path_to_kubeconfig>` environment variable configured for your target cluster

## Example

We are going to create a SealedSecret resource that can be safely stored in your deployment Chart.

First we need to create a standard Kubernetes Secret resource for the sensitive data (in this example it is a username and password pair called `creds`). By passing `--dry-run` we are instructing `kubectl` to not communicate with the remote api-server at all, and just act locally.

```
kubectl create secret generic \
  -n mynamespace creds \
  --dry-run \
  --from-literal=username=jeff \
  --from-literal=password=shhhhh \
  -o yaml > creds.yaml
```

Next we "seal" the secret using `kubeseal`

```
aws-vault exec $PROFILE -- kubeseal \
  --controller-namespace secrets-system < creds.yaml > creds-sealed.yaml
```

which will output a `SealedSecret` resource `creds-sealed.yaml` like:

```
apiVersion: bitnami.com/v1alpha1
kind: SealedSecret
metadata:
  creationTimestamp: null
  name: creds
  namespace: mynamespace
spec:
  encryptedData:
    username: 9UK7Puq03u+wLqvEeB5uUKMYWaFLyJiPF...
    password: AgC2tl+7o6hVsqGR9JInF+wu1CF9UK7Pu...
```

You can then copy this SealedSecret into your deployment Chart ready to get merged/deployed!

## One-liner

You can avoid writing the unecrypted data to an intermediate file by performing the sealing in a single step:

```
kubectl create secret generic \
  -n mynamespace creds \
  --dry-run \
  --from-literal=username=jeff \
  --from-literal=password=shhhh \
  -o yaml \
| aws-vault exec $PROFILE -- kubeseal \
  --controller-namespace secrets-system
```

Which will output the `SealedSecret` resource for a secret named `creds` with `username` and `password` key/values.

