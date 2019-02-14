# Set up monitoring and alerting

You can use Prometheus and Grafana for monitoring and alerting for your apps on the GDS Supported Platform.

### Access Prometheus or Grafana using kubectl

You must have:

- access to the AWS account that your cluster is located in
- installed the [AWS IAM Authenticator](link)
- installed [kubectl](link)

Neither Prometheus nor Grafana are publicly accessible to everyone through the internet. You must access Prometheus or Grafana through the associated Kubernetes service.

1. Run the following command to get the name of the services and ports that expose Prometheus and Grafana:

    ```
    aws vault exec <AWS_PROFILE_NAME> -- kubectl -n monitoring-system get services
    ```

    where `<AWS_PROFILE_NAME>`is your AWS profile name.

    In the following example output from this command:
    - the `monitoring-system-promethe-prometheus` service that exposes Prometheus is on port `9090`
    - the `monitoring-system-grafana` service that exposes Grafana is on port `80`

    ```
    NAME                                         TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)    AGE
    monitoring-system-grafana                    ClusterIP   10.3.0.125   <none>        80/TCP     20h
    monitoring-system-kube-state-metrics         ClusterIP   10.3.0.209   <none>        8080/TCP   20h
    monitoring-system-promethe-operator          ClusterIP   None         <none>        8080/TCP   20h
    monitoring-system-promethe-prometheus        ClusterIP   10.3.0.20    <none>        9090/TCP   20h
    monitoring-system-prometheus-node-exporter   ClusterIP   10.3.0.152   <none>        9100/TCP   20h
    prometheus-operated                          ClusterIP   None         <none>        9090/TCP   20h
    ```

1. Run the following to connect to the service that exposes either Prometheus or Grafana:

    ```
    aws-vault exec <AWS_PROFILE_NAME> -- kubectl port-forward service/<SERVICE_NAME> -n monitoring-system <PORT_1>:<PORT_2>
    ```

    where:
    - `<SERVICE_NAME>` is the name of the service that exposes either Prometheus or Grafana
    - `<PORT_1>` is a port available on your local machine
    - `<PORT_2>` is the port that exposes either Prometheus or Grafana

1. Go to `http://127.0.0.1:<PORT_1>` to see the interface.

    To access Grafana you must also sign in with:
    - username: admin
    - password: password

### Set up your app for monitoring

To set up your app for monitoring, you must:
- ensure that your app exposes metrics on the `/metrics` endpoint
- enable Prometheus to scrape metrics from that endpoint

1. You must configure your app to expose metrics on the `/metrics` endpoint. For more information on exposing metrics, refer to the Prometheus documentation on:

    - [Client libraries](https://prometheus.io/docs/instrumenting/clientlibs/)
    - [Exporters and integrations](https://prometheus.io/docs/instrumenting/exporters/) [external links]

1. To enable Prometheus to scrape metrics automatically, you must set up a [`ServiceMonitor`](https://github.com/coreos/prometheus-operator/blob/master/Documentation/user-guides/getting-started.md#related-resources) [external link] as a link between Prometheus and your app.

    To do this, you must add a `ServiceMonitor` Kubernetes resource to your deployment chart repo. A chart is a collection of files that describes a related set of Kubernetes resources.

    To create a `ServiceMonitor` resource, refer to this [example `ServiceMonitor` for the `gsp-canary` app](https://github.com/alphagov/gsp-canary-chart/blob/master/charts/gsp-canary/templates/service-monitor.yaml) [external link].

Refer to the following documentation:

- the [Helm documentation on charts](https://docs.helm.sh/developing_charts/) [external link] for more information on deployment charts
- the [`ServiceMonitor` design document](https://github.com/coreos/prometheus-operator/blob/master/Documentation/design.md#servicemonitor) for more information on this type of resource

### Create and edit alerts

Decide which alerts you want to receive by:

- reading the [GDS Way alerting principles](https://gds-way.cloudapps.digital/standards/alerting.html#alerting)
- contacting us using [re-GSP-team@digital.cabinet-office.gov.uk](mailto:re-GSP-team@digital.cabinet-office.gov.uk) or through the [Slack channel](https://gds.slack.com/messages/CDA7YSP0D/details/)

For more information on creating and editing alerts, refer to the [Prometheus Operator for Kubernetes documentation on alerting](https://github.com/coreos/prometheus-operator/blob/master/Documentation/user-guides/alerting.md) [external link].

### Further information

Refer to the following documentation for more information on:

- [how alerting works in Prometheus](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/) [external link]
- [writing an expression in PromQL](https://prometheus.io/docs/prometheus/latest/querying/basics/) for your alerting rules
