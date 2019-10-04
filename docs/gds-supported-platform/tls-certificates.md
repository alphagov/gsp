# Public TLS Certificates

All apps in GSP are protected by TLS.  There are two options for this:
using cluster-provided certificates, or providing your own.

## Using cluster-provided certificates

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

## Providing your own certificates

If you want to provide your own certificate, you will need:

 - a `SealedSecret` resource containing the certificate and key
 - a `Gateway` resource defining the ingress domain

We assume:

 - you are in a namespace called `my-namespace`
 - your certificate and key are in files called `my-cert.pem` and
   `my-cert.key`
 - you want to create a `SealedSecret` called `my-custom-cert` in a
   file called `my-custom-cert.yaml`
 - you are creating the certificate for the `sandbox` cluster
 - the certificate is for `my-custom-domain.example.com`.

To create a `SealedSecret` with your certificate, run:

    kubectl create -n my-namespace secret generic my-custom-cert --dry-run --from-file=cert=my-cert.pem --from-file=key=my-cert.key --output yaml | gds sandbox seal --format yaml > my-custom-cert.yaml

Here is an example of a `Gateway` that uses this certificate:

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: custom-certificate-gateway
  namespace: my-namespace
  annotations:
    externaldns.k8s.io/namespace: my-namespace
spec:
  selector:
    istio: my-namespace-ingressgateway
  servers:
  - port:
      number: 443
      name: https
      protocol: HTTPS
    tls:
      mode: SIMPLE
      credentialName: my-custom-cert
    hosts:
      - "my-custom-domain.example.com"
```

The line `istio: my-namespace-ingressgateway` (replace `my-namespace`
with your actual namespace name) selects the ingressgateway in your
namespace.

The line `credentialName: my-custom-cert` tells the `Gateway` to use
the `SealedSecret` you created above for TLS termination.

The line `- "my-custom-domain.example.com"` tells the `Gateway` to
listen on the virtual host corresponding to the certificate domain
name.

### Custom domains

As the example shows, if you provide your own certificates, you can
use any custom domain you choose; you do not have to use the cluster
`govsvc.uk` domain.  To use a custom domain, create a CNAME record
from your domain to the external load balancer for your namespace.

To find the external load balancer name, run:

    kubectl -n my-namespace get svc my-namespace-ingressgateway

and look under the `EXTERNAL-IP` column.

[cert-manager]: https://docs.cert-manager.io/en/latest/
[gsp-canary]: https://github.com/alphagov/gsp/tree/master/components/canary
[LetsEncrypt]: https://letsencrypt.org/
