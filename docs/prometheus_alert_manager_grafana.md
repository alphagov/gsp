# Set up monitoring and alerting

You can use either Prometheus or Grafana to monitor your apps on the GDS Supported Platform.

### Set up your app for monitoring

You must ensure that your app exposes the `/metrics` endpoint, and that this endpoint responds with compatible metrics. Refer to the following documentation for more information on compatible metrics:

- [Prometheus](https://prometheus.io/docs/instrumenting/writing_exporters/#metrics) [external link]
- [Grafana](http://docs.grafana.org/administration/metrics/) [external link]

To do this, you must create a service monitor that attaches your app (with the exposed endpoint) to your monitoring program. This enables the monitor to scrape metrics automatically.

Refer to the following `prometheus-operator` documentation for examples of this:

- [Prometheus](https://github.com/coreos/prometheus-operator/blob/master/Documentation/user-guides/running-exporters.md#generic-servicemonitor-example) [external link]
- [Grafana](https://github.com/coreos/prometheus-operator/tree/master/helm/grafana#adding-grafana-dashboards) [external link]

### Create and edit alerts

Decide which alerts you want to receive by:

- reading the [GDS Way alerting principles](https://gds-way.cloudapps.digital/standards/alerting.html#alerting)
- contacting us using [re-GSP-team@digital.cabinet-office.gov.uk](mailto:re-GSP-team@digital.cabinet-office.gov.uk) or through the [Slack channel](https://gds.slack.com/messages/CDA7YSP0D/details/)

Refer to the following documentation for information on creating and editing alerts:

- [Prometheus](https://github.com/coreos/prometheus-operator/blob/master/Documentation/user-guides/alerting.md) [external link]
- [Grafana](http://docs.grafana.org/alerting/rules) [external link]

### Access Prometheus using kubectl

You must have:

- access to the AWS account that your cluster resides on
- access to the service you are trying to monitor
- installed the [AWS IAM Authenticator](link)
- installed [kubectl](link)

Neither Prometheus nor Grafana are publicly accessible to everyone through the internet. To monitor your apps, you must tunnel to the service that exposes your monitor.

1. Run the following command to get the name of the service and port that exposes your monitor:

    ```
    kubectl -n NAMESPACE get services
    ```

1. Run the following to connect to the service that exposes your monitor:

    ```
    aws-vault exec <AWS_PROFILE_NAME> -- kubectl port-forward service/<SERVICE_NAME> -n monitoring-system <PORT_1>:<PORT_2>
    ```

    where:
    - `<AWS_PROFILE_NAME>`is your AWS profile name
    - `<SERVICE_NAME>` is the name of the service that exposes your monitor
    - `<PORT_1>` is the port exposed by the service
    - `<PORT_2>` is a port available on your local machine

1. Go to `http://127.0.0.1:<PORT_2>` to see the monitor interface.


### Further information

Refer to the following documentation for more information on:

- [how alerting works in Prometheus](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/) [external link]
- [writing an expression in PromQL](https://prometheus.io/docs/prometheus/latest/querying/basics/) for your alerting rules.
- what else?
