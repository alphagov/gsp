# Public TLS Certificates

You can provision TLS certificates using [cert-manager][] in GSP.

By default, the GSP has a `ClusterIssuer` named `letsencrypt-r53` that is configured to provision TLS certificates supplied by [LetsEncrypt][] via the DNS01 ACME challenge. For example, to add a TLS certificate for the [gsp-canary][] in the sandbox cluster use the following kube yaml:

```yaml
apiVersion: certmanager.k8s.io/v1alpha1
kind: Certificate
metadata:
  name: sandbox-gsp-canary-ingress
  namespace: sandbox-main
spec:
  acme:
    config:
    - dns01:
        provider: route53
      domains:
      - canary.london.sandbox.govsvc.uk
  dnsNames:
  - canary.london.sandbox.govsvc.uk
  issuerRef:
    kind: ClusterIssuer
    name: letsencrypt-r53
  secretName: sandbox-gsp-canary-ingress-certificate
```

> **Note:** cert-manager will need to be able to modify the DNS of the domains listed in the certificate in order to perform the DNS challenge. At the time of writing that only applies to the cluster domain.

[cert-manager]: https://docs.cert-manager.io/en/latest/
[gsp-canary]: https://github.com/alphagov/gsp/tree/master/components/canary
[LetsEncrypt]: https://letsencrypt.org/
