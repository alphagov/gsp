# Public DNS

The GSP allows applications to set DNS records (route53) based on istio `Gateway` resources on a per-namespace basis. The `Gateway` resource needs to carry the correct annotation (with key `externaldns.k8s.io/namespace`) to ensure the `external-dns` instance adds the corresponding A record(s).

To set the DNS entry for the [gsp-canary][] in the sandbox cluster to `canary.london.sandbox.govsvc.uk` use the following kube yaml:

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  annotations:
    externaldns.k8s.io/namespace: sandbox-main
  name: sandbox-gsp-canary-ingress
  namespace: sandbox-main
spec:
  selector:
    istio: sandbox-main-ingressgateway
  servers:
  - hosts:
    - canary.london.sandbox.govsvc.uk
    port:
      name: https
      number: 443
      protocol: HTTPS
    tls:
      credentialName: sandbox-gsp-canary-ingress-certificate
      mode: SIMPLE
      privateKey: sds
      serverCertificate: sds
```

For more details on the features of `external-dns` see the [external-dns] documentation.

Each namespace has an instance of [external-dns][] running that will configure DNS A records to point at the load balancer created for the ingressgateway in the same namespace. The `external-dns` instance in each namespace is limited to only create records in the route53 zone for which the cluster is authoritative (`london.sandbox.govsvc.uk` in the above example).

The locations of the TLS certificates will depend on the installation context. This example above uses istio's [secret discovery service][] to dynamically load the certificate from a secret with the name given in `credentialName`.

[external-dns]: https://github.com/kubernetes-incubator/external-dns
[gsp-canary]: https://github.com/alphagov/gsp/tree/master/components/canary
[secret discovery service]: https://istio.io/docs/tasks/traffic-management/ingress/secure-ingress-sds/
