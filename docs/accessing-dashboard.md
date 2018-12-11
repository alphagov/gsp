# Accessing the Kubernetes Dashboard

Each cluster will have a deployment of the kubernetes dashboard available. This will not have an ingress configured and so will not be accessible from outside the cluster. It will however work via a proxy tunnel. The service account the dashboard is running under has very limited access (making the dashboard useless) so authenticating against the dashboard is also required. This guide will walk through how to connect to the dashboard and authenticate using your IAM credentials.

## Requirements

* [aws-iam-authenticator](https://github.com/kubernetes-sigs/aws-iam-authenticator)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* a kubeconfig as set up by the cluster terraform

### Suggested

* [aws-vault](https://github.com/99designs/aws-vault)
* jq

## Steps

1. Set up the tunnel

        export KUBECONFIG=<path_to_kubeconfig>
        kubectl proxy
1. Load the dashboard  
In a browser, go to: http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/#!/login
1. Get a token  
In another terminal:

        aws-vault exec run-sandbox -- aws-iam-authenticator token -i <cluster_id> | jq -r .status.token | pbcopy
1. Login  
Select "Token" and paste the token that was copied in the previous step into the box. Click "Sign In".
