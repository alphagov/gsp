# OPA policies

## Running Tests

All together:

```
$ opa test policies
PASS: 21/21
```

Or, individually:

```
$ opa test policies/digests-on-images
PASS: 5/5
```
```
$ opa test policies/restrict-special-nodes
PASS: 11/11
```
```
$ opa test policies/isolate-tenant-istio-resources
PASS: 5/5
```
