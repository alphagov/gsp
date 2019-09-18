# Grafana inside GSP

## Where to find Grafana

Browse to https://grafana.london.{name of cluster}.govsvc.uk/

## Login via Google

When your cluster is created you should provide the deployer pipeline the Google OAuth client ID and secret in `google-oauth-client-id` and `google-oauth-client-secret`.

## Getting the default admin password of a new GSP Grafana instance

Using gds-cli, as an admin in the kubernetes cluster:
gds {name of cluster} kubectl get -n gsp-system secret gsp-grafana -o json | jq -r '.data["admin-password"]' | base64 -D -

