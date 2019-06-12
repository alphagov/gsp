# Deployer pipeline

Example for pushing a deployer pipeline to provision a sandbox cluster can be found in `./hack/set-deployer-pipeline.sh`:

```
CLUSTER_CONFIG=./pipelines/examples/clusters/sandbox.yaml USER_CONFIGS=./pipelines/examples/users/ ./hack/set-deployer-pipeline.sh
```
