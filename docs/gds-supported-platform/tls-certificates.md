# Public TLS Certificates

Your GSP app will need a TLS certificate in order to serve HTTPS
traffic.  There are two options:

 - use a cluster-provided certificate
 - provide your own certificate

Cluster-provided certificates require less effort because you don't
need to provision a certificate for yourself.  However, the cluster
cannot at this point provide certificates for custom domains (that is,
non-`.govsvc.uk` domains).  If you wish to use a custom domain, you
must provide your own TLS certificate.

## Using cluster-provided certificates

You can provision TLS certificates using [cert-manager][] in GSP.

By default, the GSP has a `ClusterIssuer` named `letsencrypt-r53` that is configured to provision TLS certificates supplied by [LetsEncrypt][] via the DNS01 ACME challenge. For example, to add a TLS certificate for the [gsp-canary][] in the sandbox cluster use the following kube yaml:

```yaml
apiVersion: cert-manager.io/v1alpha2
kind: Certificate
metadata:
  name: sandbox-gsp-canary-ingress
  namespace: sandbox-main
spec:
  dnsNames:
  - canary.london.sandbox.govsvc.uk
  issuerRef:
    kind: ClusterIssuer
    name: letsencrypt-r53
  secretName: sandbox-gsp-canary-ingress-certificate
```

> **Note:** cert-manager will need to be able to modify the DNS of the domains listed in the certificate in order to perform the DNS challenge. At the time of writing that only applies to the cluster domain.

## Using a custom domain by providing your own certificate

If you want to use a custom domain, you must provide your own
certificate.  To do this, you must create a:

 - `SealedSecret` resource with the certificate and key
 - `Gateway` resource to listen on the domain
 - `CNAME` record from your custom domain to the namespace ingressgateway

### Create a `SealedSecret` with the certificate and key

To create a `SealedSecret` with your certificate, run:

    kubectl create -n <NAMESPACE> secret generic <CERTNAME> --dry-run --from-file=cert=<CERTFILE> --from-file=key=<KEYFILE> --output yaml | gds <CLUSTER> seal --format yaml

Where:

 - `<NAMESPACE>` is your namespace
 - `<CERTNAME>` is the name you will give to your `SealedSecret`
 - `<CERTFILE>` and `<KEYFILE>` are the filenames of your certificate
   and key in PEM format
 - `<CLUSTER>` is the GSP cluster you are targeting (for example,
   `verify`)

### Create a `Gateway` to listen on the domain

To use this `SealedSecret`, create a `Gateway` with the following
yaml:

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: <NAME>
  namespace: <NAMESPACE>
  annotations:
    externaldns.k8s.io/namespace: <NAMESPACE>
spec:
  selector:
    istio: <NAMESPACE>-ingressgateway
  servers:
  - port:
      number: 443
      name: https
      protocol: HTTPS
    tls:
      mode: SIMPLE
      credentialName: <CERTNAME>
    hosts:
      - "<CUSTOM_DOMAIN>"
```

Where:

 - `<NAME>` is the name you will give this `Gateway` resource
 - `<CUSTOM_DOMAIN>` is the fully-qualified domain name for your
   custom domain for your certificate (for example,
   `my-custom-domain.example.com`)

Note: the line `istio: <NAMESPACE>-ingressgateway` selects the
ingressgateway in your namespace.

### Create a CNAME record from your custom domain to the namespace ingressgateway

In the DNS configuration for your custom domain, create a CNAME record
from your custom domain to `<NAMESPACE>.london.<CLUSTER>.govsvc.uk`.

[cert-manager]: https://docs.cert-manager.io/en/latest/
[gsp-canary]: https://github.com/alphagov/gsp/tree/master/components/canary
[LetsEncrypt]: https://letsencrypt.org/
