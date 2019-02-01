# Troubleshooting app errors

When troubleshooting your app, you should use [Kubectl](https://kubernetes.io/docs/reference/kubectl/kubectl/) [external link] to access your Kubernetes cluster and find out more information.

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

### Set up a kubeconfig file

Kubeconfig files contain configuration information for kubectl, including your cluster's API endpoint, the user profile, and the public key.

You can use one kubeconfig file for multiple clusters, or one kubeconfig file per cluster. It is your choice as to how you set up your kubeconfig file(s).

The GDS Supported Platform team will provide you with a kubeconfig file after assigning you to your account. Contact us by emailing [re-GSP-team@digital.cabinet-office.gov.uk](mailto:re-GSP-team@digital.cabinet-office.gov.uk) or on [Slack](https://gds.slack.com/messages/CDA7YSP0D) if you have any questions.

Refer to the [Kubernetes documentation on accessing clusters with kubeconfig files](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/) for more information.

### Install the AWS IAM Authenticator

The AWS IAM Authenticator authenticates users for our Kubernetes clusters.

1. You must install the AWS IAM Authenticator on your local machine so kubectl can authenticate your access to a cluster.

    Refer to the [AWS IAM authenticator installation documentation](https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html) for more information.

1. Once you have installed the authenticator, sign into the authenticated AWS IAM session to access your cluster.  You can do this in multiple ways, for example by using [`aws-vault`](https://github.com/99designs/aws-vault) [external link].

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

## Further information

For more information, refer to the following documentation:

- how to [install kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- configuring [access to multiple clusters](https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/)
- the kubectl [cheat sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)
- the kubectl [reference documentation](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands) for information on other commands that you can run [external links]
