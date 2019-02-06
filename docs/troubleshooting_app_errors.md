# Troubleshooting app errors

When troubleshooting your app, you should use [kubectl](https://kubernetes.io/docs/reference/kubectl/kubectl/) [external link] to access your Kubernetes cluster and find out more information.

You should use the [Kubernetes dashboard](https://github.com/kubernetes/dashboard) [external link] to manage your cluster, the apps running in your cluster, and also help troubleshoot those apps.

You can also access your Kubernetes cluster using [Cloudwatch](https://aws.amazon.com/cloudwatch/) [external link] or [Logit](https://logit.io/) [external link].

## Access your cluster using kubectl

To use kubectl, you must:

- have a [GDS users AWS account](https://reliability-engineering.cloudapps.digital/iaas.html#amazon-web-services-aws)
- install kubectl
- set up a kubeconfig file for the cluster you need to access
- install the AWS IAM Authenticator

### Install kubectl

For MacOS, run the following command to install kubectl:

```
brew install kubernetes-cli
```

Refer to the [Kubernetes documentation on installing kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) [external link] for other operating systems.

### Install the AWS IAM Authenticator

The AWS IAM Authenticator authenticates users for our Kubernetes clusters.

1. You must install the AWS IAM Authenticator on your local machine so kubectl can authenticate your access to a cluster.

    Refer to the [AWS IAM authenticator installation documentation](https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html) for more information.

### Sign into an authenticated AWS IAM session

Once you have installed the authenticator, you must sign into the authenticated AWS IAM session to access your cluster. You should use [`aws-vault`](https://github.com/99designs/aws-vault) [external link] to do this.

1. Install `aws-vault` on macOS by running:

    ```
    brew cask install aws-vault
    ```

    Refer to the [`aws-vault` installation documentation](https://github.com/99designs/aws-vault#installing) [external link] for other ways to install this program.

1. Add your credentials to `aws-vault` by running:

    ```
    aws-vault add gds-users
    ```

1. Add your Access Key ID and your Secret Access Key to a new keychain. You can find `aws_access_key_id` and `aws_secret_access_key` by signing into the GDS AWS account using the AWS console.

    You will see the message `Added credentials to profile "gds-users" in vault` in your command line when this is complete.

Refer to the [Reliability Engineering documentation on accessing AWS accounts](https://reliability-engineering.cloudapps.digital/iaas.html#access-aws-accounts) for more information.

### Set up a kubeconfig file

Kubeconfig files contain configuration information for kubectl, including your cluster's API endpoint, the user profile, and the public key.

The GDS Supported Platform team will provide you with a secure kubeconfig file after assigning you to your account. This kubeconfig file will only work if you are signed into an authenticated AWS account.

Contact us at [re-GSP-team@digital.cabinet-office.gov.uk](mailto:re-GSP-team@digital.cabinet-office.gov.uk) if you have any questions.

Refer to the [Kubernetes documentation on accessing clusters with kubeconfig files](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/) for more information.

## Use kubectl to get information on your cluster

When you are signed into an authenticated AWS IAM session, you can run different commands to get information on your cluster. The most common commands are summarised in this section.

These commands are specified by [namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/) and [pod](https://kubernetes.io/docs/concepts/workloads/pods/pod/) [external links].

### Get information about a namespace

Run the following in the command line to get a summary of information about the specified namespace:

```
kubectl -n NAMESPACE get all
```

### Get information about pods within a namespace

Run the following in the command line to list all pods within a namespace:

```
kubectl -n NAMESPACE get pods
```

Run the following to get information about a specific pod within a namespace:

```
kubectl -n NAMESPACE get pods PODNAME
```

Run the following to get information about multiple pods within a namespace:

```
kubectl -n NAMESPACE get pods PODNAME_1:PODNAME_2:...:PODNAME_N
```

### Get information about events for a namespace

Run the following in the command line to get information for all events for a namespace:

```
kubectl -n NAMESPACE get events
```

Run the following to get detailed information on a pod:

```
kubectl -n NAMESPACE describe pod PODNAME
```

### Get logs for a pod within a namespace

Run the following in the command line to get all logs for a pod within a namespace at that point in time:

```
kubectl -n NAMESPACE logs PODNAME
```

This command picks up logs that are written to `STDOUT` and `STDERR`. If you write your logs to another destination or file, this command will not pick up those logs.

Run the following to get a continually updating view of all logs for a pod within a namespace:

```
kubectl -n NAMESPACE logs PODNAME --follow
```

#### Use Kail to get logs

You can also use [Kail](https://github.com/boz/kail) [external link] to get logs by namespace, pod, ingress or container.

1. Install Kail by running the following in the command line:

    ```
    brew tap boz/repo
    brew install boz/repo/kail
    ```

1. Set up an alias by running:

    ```
    alias kail='aws-vault exec run-sandbox -- /usr/local/bin/kail'
    ```
    _QP: Why would you do this?_

1. Run the following to get logs by namespace, pod, ingress or container:

    ```
    xyz
    ```

## Manage your cluster and apps using the Kubernetes dashboard

Each cluster has a Kubernetes dashboard deployment available. The dashboard does not have an ingress configured and so is not accessible from outside the cluster.

You must connect to the dashboard using a proxy tunnel, and authenticate against the dashboard using your AWS IAM credentials.

Before you start, you must have:

- installed kubectl
- a configured kubeconfig file
- installed the AWS IAM authenticator and signed into a authenticated AWS IAM session

We recommend that you use [`aws-vault`](https://github.com/99designs/aws-vault#installing) [external link] to sign into the authenticated AWS IAM session to access your cluster.

1. Set up the proxy tunnel by running the following in the command line:

    ```
    export KUBECONFIG=<path_to_kubeconfig>
    kubectl proxy
    ```
    _QP: Is the `path_to_kubeconfig` a placeholder?_

1. Get an authentication token by running:

    ```
    aws-vault exec run-sandbox -- aws-iam-authenticator token -i <cluster_id> | jq -r .status.token | pbcopy
    ```
    Copy the outputted token.

1. Open your web broswer and go to [http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/#!/login](http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/#!/login) [external link] to load the dashboard.

1. Select __Token__ and paste the copied token into the box.

    _QP: What does the screen look like?_

1. Select __Sign In__ to log into the dashboard.




## Further information

For more information, refer to the following documentation:

- how to [install kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- configuring [access to multiple clusters](https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/)
- the kubectl [cheat sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)
- the kubectl [reference documentation](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands) for information on other commands that you can run [external links]
