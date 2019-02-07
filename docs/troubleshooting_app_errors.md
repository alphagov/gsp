# Troubleshooting app errors

When troubleshooting your app, you should use [kubectl](https://kubernetes.io/docs/reference/kubectl/kubectl/) [external link] to access your Kubernetes cluster to find out more information.

You can use the [Kubernetes dashboard](https://github.com/kubernetes/dashboard) [external link] instead to access your Kubernetes cluster using a GUI rather than the command line.

## Access your cluster using kubectl

To use kubectl, you must:

- have a [GDS users AWS account](https://reliability-engineering.cloudapps.digital/iaas.html#amazon-web-services-aws)
- install kubectl
- install the AWS IAM Authenticator
- have a kubeconfig file for the cluster you need to access

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

### Get a kubeconfig file

Kubeconfig files contain configuration information for kubectl, including your cluster's API endpoint, the user profile, and the public key.

The GDS Supported Platform team will provide you with a secure kubeconfig file after assigning you to your account. This kubeconfig file will only work if you are signed into an authenticated AWS account.

Contact us at [re-GSP-team@digital.cabinet-office.gov.uk](mailto:re-GSP-team@digital.cabinet-office.gov.uk) if you have any questions.

Refer to the [Kubernetes documentation on accessing clusters with kubeconfig files](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/) for more information.

## Use kubectl to get information on your cluster

When you are signed into an authenticated AWS IAM session, you can run different commands to get information on your cluster. The most common commands are summarised in this section.

These commands are specified by [namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/) and [pod](https://kubernetes.io/docs/concepts/workloads/pods/pod/) [external links].

Before you can run these commands, you must run the following:

```
export KUBECONFIG=<PATH_TO_KUBECONFIG_FILE>
```

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

This command only gets logs from one source at a time.

### Get logs from multiple sources at the same time

You can use the [Kail plugin](https://github.com/boz/kail) [external link] to get logs from multiple sources at the same time.

1. Run the following in the command line to tell Homebrew where the Kail code is located:

    ```
    brew tap boz/repo
    ```

1. Run the following to install Kail:

    ```
    brew install boz/repo/kail
    ```

1. Run the following to make Kail ready to receive log parameters:

    ```
    aws-vault exec run-sandbox -- /usr/local/bin/kail
    ```

    Refer to the [Kail documentation](https://github.com/boz/kail/blob/master/README.md) for more information on the log parameters.

## Access your cluster using the Kubernetes dashboard

Each cluster has a Kubernetes dashboard deployment available. You can use the dashboard find out information about the cluster using a GUI

You cannot access the dashboard from outside the cluster. You must connect to the dashboard using a proxy tunnel, and authenticate against the dashboard using your AWS IAM credentials.

Before you start, you must have:

- installed kubectl
- a configured kubeconfig file
- installed the AWS IAM authenticator and signed into a authenticated AWS IAM session

We recommend that you use [`aws-vault`](https://github.com/99designs/aws-vault#installing) [external link] to sign into the authenticated AWS IAM session to access your cluster.

1. Set up the proxy tunnel by running the following in the command line:

    ```
    export KUBECONFIG=<PATH_TO_KUBECONFIG_FILE>
    kubectl proxy
    ```

    where `KUBECONFIG_FILE` is the location of your kubeconfig file.

1. Get an authentication token by running:

    ```
    aws-vault exec run-sandbox -- aws-iam-authenticator token -i <CLUSTER_ID> | jq -r .status.token | pbcopy
    ```

    The `CLUSTER_ID` is in the kubeconfig file. In the following kubeconfig file example, the cluster ID is `johnsmith.run-sandbox.aws.ext.govsandbox.uk`.

    ```
    users:
     - name: johnsmith-user
       user:
         exec:
           apiVersion: client.authentication.k8s.io/v1alpha1
           command: aws-iam-authenticator
           args:
             - "token"
             - "-i"
             - "johnsmith.run-sandbox.aws.ext.govsandbox.uk"
    ```

1. Copy the outputted token.

1. Open your web browser and go to [http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/#!/login](http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/#!/login) [external link] to load the dashboard.

1. Select __Token__ and paste the copied token into the box.

1. Select __Sign In__ to log into the dashboard.

You can now view your cluster and the apps running in that cluster, and find out information to help troubleshoot those apps.

## Further information

For more information, refer to the following documentation:

- how to [install kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- configuring [access to multiple clusters](https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/)
- the kubectl [cheat sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)
- the kubectl [reference documentation](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands) for information on other commands that you can run [external links]
