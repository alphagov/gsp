# Troubleshooting app errors

Notes from scoping session:

- How to access my cluster using Kubectl
- Learn how to use Kubectl effectively
- Repurpose some of the diagnostics and debugging firebreak content

## How to access my cluster using Kubectl

when troubleshooting your app, accessing your cluster using Kubectl is a standard first step

Kubectl is a standard tool for interacting with kubernetes, so there is existing documentation

to use Kubectl you need to:

1. install Kubectl binary
1. install AWS IAM authenticator
1. create Kubectl kubeconfig file that is set up for the cluster (either locally or remotely TBC)

you can use one kubeconfig file to multiple clusters, or one kubeconfig file for one cluster. it's up to the user, we don't mandate - we might make a recommendation but don't know yet.

kubeconfig files contain API endpoint for cluster, user profile, public key and other information.

### AWS IAM authenticator

AWS IAM is our authentication method for our clusters.

assumption - every member of the service team who needs to access kubernetes has been given access to kubernetes, we'll need to know their IAM details (pre-requisite)

in order for kubectl to use IAM to authenticate a user to perform any action on the cluster, the user also needs AWS IAM authenticator installed on their local machine.

prerequisite: kubectl installed (macOS `brew install kubernetes-cli`)

prerequisite: kubeconfig file provided to user __QP: how will team provide users with kubeconfig files? TBC__

to install AWS IAM authenticator, go to https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html

when you run kubectl, you need to be in an authenticated aws iam session. there are mutliple ways to achieve this, for example using `aws-vault` - https://github.com/99designs/aws-vault.

### namespace

You must always specify the namespace when you run a command `kubectl -n NAMESPACE`

### get description of everything

`kubectl -n NAMESPACE get all`

Summary get information on all objects within a namespace for example pods, services,

### Get pods

A pod is a thing that executes.

List the pods in a namespace: `kubectl -n NAMESPACE get pods`

Get information about a specific pod: `kubectl -n NAMESPACE get pods PODNAME`

You can get information on multiple pods: `kubectl -n NAMESPACE get pods PODNAME1 PODNAME_N`

### Events

You can run `kubectl -n NAMESPACE describe pod PODNAME` to detailed information on a pod, e.g. for ports and probes. It also includes events, which might be useful for troubleshooting.

You can get detailed event information for a namespace: `kubectl -n NAMESPACE get events` - information overload

### Logs

To get logs for a pod within a namespace: `kubectl -n NAMESPACE logs PODNAME`

Note that this command will get the logs that are being written to `STDOUT` and `STDERR`. If you are writing your logs to another destination or file, this command will not pick those logs up.

To get logs for a pod within a namespace and for this to get an updating window as more logs come in: `kubectl -n NAMESPACE logs PODNAME --follow`

### equivalents

equivalent gets and describes and other for ingress, services, deployments and so on

https://kubernetes.io/docs/reference/kubectl/cheatsheet/

https://kubernetes.io/docs/reference/kubectl/cheatsheet/#viewing-finding-resources

might be useful, take a look

## Further information

How to install Kubectl - https://kubernetes.io/docs/tasks/tools/install-kubectl/

Configure Access to Multiple Clusters - https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/
