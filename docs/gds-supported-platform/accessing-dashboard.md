# Accessing the Kubernetes Dashboard

## Using the GDS CLI

1. [Install the GDS CLI](https://github.com/alphagov/gds-cli/#installation)
1. `gds sandbox dashboard`
1. Copy and paste the provided token into the dashboard login form

## Without using the GDS CLI

1. `aws-vault exec sandbox -- kubectl port-forward --namespace kube-system svc/kubernetes-dashboard 8443:443`
1. `open https://127.0.0.1:8443`
1. `aws-vault exec sandbox -- aws eks get-token --cluster-name sandbox  | jq -r .status.token`
1. Copy and paste the provided token into the dashboard login form
