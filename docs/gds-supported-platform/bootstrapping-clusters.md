# Bootstrapping GSP on-demand clusters

## Process to create a cluster

* Write a cluster-config file and values file into a branch of tech-ops-cluster-config.git, like [tech-ops-cluster-config #268](https://github.com/alphagov/tech-ops-cluster-config/pull/268/files)
  * Set `config-version` to the branch name you chose.
  * Set `config-path` to the path to your new cluster YAML file, relative to the root of tech-ops-cluster-config.
  * Set `config-values-path` to the path of your new cluster values YAML file, relative to the root of tech-ops-cluster-config.
  * Set `splunk-enabled` to "0" if you don't require Splunk in your cluster.

  * Set the following vars:
      * account-id: "011571571136"
      * account-name: "sandbox"
      * account-role-arn: "arn:aws:iam::011571571136:role/deployer"
      * github-client-id: "NOTID"
      * github-client-secret: "NOTASECRET"
      * hsm-splunk-index: "NOTAURL"
      * hsm-splunk-hec-token: "NOTATOKEN"
      * k8s-splunk-index: "NOTAURL"
      * k8s-splunk-hec-token: "NOTATOKEN"
      * google-oauth-client-secret: "NOTID"
      * google-oauth-client-id: "NOTASECRET"

* Commit and push that, open a draft PR. You don't need to get it merged.
* Ensure your fly config has a `cd-gsp` target pointing to the `gsp` team in Big Concourse.
* Change directory to top level of `gsp` repo
* `CLUSTER_CONFIG=../tech-ops-cluster-config/clusters/name-of-my-cluster.yaml ./hack/set-deployer-pipeline.sh`
* Go into Big Concourse and run the update job for the new `name-of-my-cluster-deployer` pipeline.

## Maintenance

* Remember to run the destroy pipeline (deployer pipeline in Big Concourse, `destroy` group) before you go home. This can take 20 minutes. If you have an RDS cluster (e.g. Kubernetes type `Postgres`), you may want to get rid of that too.
* You'll get releases just like sandbox does. After starting the destroy pipeline, remember to `fly -t cd-gsp pause-job -j name-of-my-cluster-deployer/deploy` so that GSP PRs merged before you return to work don't trigger your cluster to be re-deployed.

## Limitations

* To work with the cluster in the normal way you'll need gds-cli version 1.27.0.
* Service-operator'd resources will not be deleted on cluster destroy.
* After you destroy and re-create a cluster, you must `gds sandbox -c name-of-my-cluster update-kubeconfig`, or you will get errors like `Unable to connect to the server: dial tcp: lookup 8D8F3A460045AFA69F63F44F8DAB3F68.yl4.eu-west-2.eks.amazonaws.com: no such host` when trying to use kubectl.
* You can choose a custom branch to deploy Terraform/Helm with `platform-version` (and also setting `platform-resource-type` to `git` and `platform-tag-filter` to the empty string) but this does not enable the testing of e.g. the components, the differences between the deployment pipelines, or the release pipeline.
* Will not have external secrets - e.g. GitHub and Google integration. Therefore no Google login for Grafana, or GitHub login for Concourse. (NOTE: You need to be 'admin' in repo: gds-trusted-developers for these commands to function).
  * You can log into Concourse with username `pipeline-operator` and password coming from `gds sandbox -c name-of-my-cluster kubectl get -n gsp-system secret gsp-pipeline-operator -o json | jq -r '.data.concourse_password' | base64 --decode -`
  * You can log into Grafana with username `admin` and password coming from `gds sandbox -c name-of-my-cluster kubectl get -n gsp-system secret gsp-grafana -o json | jq -r '.data["admin-password"]' | base64 --decode -`

## Troubleshooting

You may run into one of the following problems:
* Ingress for gsp-system applications (e.g. Little Concourse, Grafana) refuses connections
* Ingress for your-cluster-name-main applications (e.g. Canary) refuses connections
* kubectl apply fails (e.g. inside the `cd-smoke-test` pipeline in Little Concourse) due to a control plane failure, control plane logs show it failed to connect to the `v1beta1.metrics.k8s.io` `APIService` in `kube-system`.

In these cases, try deleting all pods from the namespace in question (e.g. `gds sandbox -c name-of-my-cluster kubectl delete -n my-namespace pod --all`). Stuff should be appearing in the logs of the namespace's ingressgateway pods' istio-proxy containers if it is receiving requests.

* Clusters often come up in states with broken ingress, deleting the istio-ingressgateway pods usually fixes this, e.g. the following for Concourse/Grafana etc.: `gds sandbox -c name-of-my-cluster kubectl -n gsp-system delete pods -l app=istio-ingressgateway` - for the canary use the namespace named dynamically based on the name of your cluster followed by `-main`.
* Clusters often come up in states with broken Concourse pipeline creation, deleting the pipeline-operator pod usually fixes this: `gds sandbox -c name-of-my-cluster kubectl -n gsp-system delete pod gsp-pipeline-operator-0`

## Process to delete a cluster

* Delete any Service Operator resources in the cluster (e.g. types `Principal`, `SQS`, `S3Bucket`, `Postgres`), check that the CloudFormation stacks for those get deleted.
* Run the `destroy` job (in its own group) in the Big Concourse name-of-my-cluster-deployer pipeline
* `fly -t cd-gsp destroy-pipeline -p name-of-my-cluster-deployer`
