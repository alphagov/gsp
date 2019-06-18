# gsp-canary

The canary is a monitoring tool which is continuously changed, built, and
deployed to a Kubernetes cluster.

The canary exposes a health check endpoint with metrics that can be gathered by
Prometheus.

The intention is to smoke out any problems with the pipeline.

It should monitor itself wherever possible, including testing whether or not its
own age has passed a threshold, which might indicate that a problem exists with
the deployment process.
