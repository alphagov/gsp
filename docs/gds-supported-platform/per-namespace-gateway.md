# Per-namespace istio gateways

Each namespace has an istio `ingressgateway` to facilitate greater control and flexibility for tenant applications (for more details see the [ADR](https://github.com/alphagov/gsp/blob/master/docs/architecture/adr/ADR037-per-namespace-gateways.md)). The `ingressgateway` selector adheres to the following helm template pattern:

```yaml
{{ .Release.Namespace }}-ingressgateway
```

For example, to add a `Gateway` for the [gsp-canary][] in the `sandbox-main` namespace in the sandbox cluster use the following kube yaml (e.g. which could be rendered from a helm template):

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
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

For more details on working with istio's ingressgateways see the [ingress gateway][] docs.

By default the ingressgateway listens on TCP port 80 (for HTTP ingress traffic) and TCP port 443 (for HTTPS ingress traffic). Additional ports can be included via `*-cluster-config` values. For example, to add TCP port 3306 (for MySQL) to the ingressgateway in the `sandbox-connector-node-metadata` namespace the `values.yaml` needs to include:

```yaml
namespaces:
- name: sandbox-connector-node-metadata
  ingress:
    ports:
    - port: 3306
      name: tcp-mysql
      targetPort: 3306
```

## Migration to per-namespace gateways

During the interim period between using a central ingressgateway in the `istio-system` namespace and having an ingressgateway in all tenant namespaces the per-namespace gateways will be opt-in. For example, to opt-in the `sandbox-connector-node-metadata` namespace the `values.yaml` for the cluster (in `*-cluster-config`) needs to include:

```yaml
namespaces:
- name: sandbox-connector-node-metadata
  ingress:
    enabled: true
```

[gsp-canary]: https://github.com/alphagov/gsp/tree/master/components/canary
[ingress gateway]: https://istio.io/docs/tasks/traffic-management/ingress/ingress-control/
