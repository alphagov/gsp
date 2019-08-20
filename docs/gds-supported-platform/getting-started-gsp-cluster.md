# Getting started with a GDS Supported Platform cluster

The GDS Supported Platform (GSP) is a platform to host Government Digital Service (GDS) services. The GSP is based on [Docker](https://docs.docker.com/) and [Kubernetes](https://kubernetes.io/).

These instructions explain how to set up a remote instance of the GSP to host an app.

## Before you start

You should have:

- a [Docker image](https://docs.docker.com/engine/reference/commandline/images/) of your app built in line with the [12 factor principles](https://docs.cloud.service.gov.uk/architecture.html#12-factor-application-principles)
- access to a [Kubernetes cluster](https://github.com/alphagov/gsp/blob/master/docs/gds-supported-platform/troubleshooting_app_errors.md) created by Tech Ops

_TBC: deploying multiple apps on the same cluster that are connected to each other - separate story on how to handle networking_

## Create a GitHub repository

You must create a repository on GitHub to store your Deployment configuration. The GDS Supported Platform only works with repositories stored on GitHub.

## Request a namespace and cluster configuration

Your service team's Site Reliability Engineer (SRE) must [ask GDS Tech Ops](re-GSP-team@digital.cabinet-office.gov.uk) for a new namespace. You must provide the following information:

* your [GitHub repository address](https://help.github.com/en/articles/about-remote-repositories)
* the name of the [namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/) you want to create
* the email addresses, GitHub usernames, and [GPG public keys](https://www.gnupg.org/gph/en/manual/c14.html) (if applicable) of the users who should get read-only access to the namespace

GDS Tech Ops will create a namespace and a cluster configuration, or `clusterconfig`. The `clusterconfig` contains the following code and default values:

```
- name: CLUSTER_CONFIG-TEAM_NAME-CONCOURSE_BUILD_PIPELINE
  owner: alphagov
  repository: REPO_URL
  path: ci/build
```
You must change the values in the `clusterconfig` to configure the Kubernetes cluster.

## Create a Concourse build pipeline

A Concourse build pipeline enables continuous integration and development, and deploys your app.

1. Create a root directory in your GitHub repository.

1. Create a `ci/build` directory in the new root directory.

1. Create a Concourse build pipeline `.yaml` file in the `ci/build` directory. The `.yaml` file should:

    - define the `github_source`, `harbor_source` and Concourse `resource_types`
    - create the Docker container resource
    - create the `build job` that builds all containers and does unit testing
    - define the `packaging`
    - generate chart values and the manifest file
    - create the release job

1. Create a deployment pipeline `.yaml` file in the `ci/build` directory. The deployment pipeline `.yaml` file uses the artifacts created by the Concourse build pipeline.

Refer to the following Concourse documentation for more information on creating Concourse build pipelines:

- [pipelines](https://concourse-ci.org/pipelines.html)
- [resources](https://concourse-ci.org/resources.html) and [resource types](https://concourse-ci.org/resource-types.html)
- [jobs](https://concourse-ci.org/jobs.html)
- [tasks](https://concourse-ci.org/tasks.html)
- [builds](https://concourse-ci.org/builds.html)

## Create a Helm chart

Kubernetes resources describe the configuration of the app that you are running. You define these resources as `.yaml` files and collect them together in a packaging format called a [Helm chart](https://helm.sh/docs/developing_charts/).

You create Helm charts as files in a directory with the following structure:

```
chart/
├── Chart.yaml
└── templates/
    ├── deployment.yaml
    ├── destinationrule.yaml
    ├── service.yaml
    └── virtualservice.yaml
```

### Create the Helm chart directory and yaml files

1. Create a `chart` directory inside the root directory of the GitHub repository.

1. Create a `Chart.yaml` file in the `chart` directory with the following:

    ```
    apiVersion: v1
    appVersion: "1.0"
    description: CHART_DESCRIPTION
    name: CHART_NAME
    version: 0.1.0
    ```

    This file defines metadata about the chart.

1. Create a `templates` directory in the `chart` directory to store the Kubernetes object definitions.

1. Create a `values.yaml` file in the root directory to set the default values for your desired chart variables.

### Create a Kubernetes Deployment resource

You run an app in the local GSP instance by creating a [Kubernetes Deployment resource object](https://kubernetes.io/docs/concepts/#kubernetes-objects). This object defines your app and its routes, databases and all other relevant information. You describe a Kubernetes Deployment in a `.yaml` file.

1. Create a `deployment.yaml` file in the `templates` directory. The following example uses an [nginx](https://hub.docker.com/_/nginx/) container image called `myapp`. Replace this nginx container image with your app image:

    ```
    apiVersion: apps/v1beta2
    kind: Deployment
    metadata:
      name: {{ .Release.Name }}-myapp
      labels:
        app.kubernetes.io/name: myapp
        app.kubernetes.io/instance: {{ .Release.Name }}
    spec:
      replicas: {{ .Values.replicas }}
      selector:
        matchLabels:
          app.kubernetes.io/name: myapp
          app.kubernetes.io/instance: {{ .Release.Name }}
      template:
        metadata:
          labels:
            app.kubernetes.io/name: myapp
            app.kubernetes.io/instance: {{ .Release.Name }}
        spec:
          containers:
            - name: myappcontainer
              image: "nginx:latest" #Replace this with your app image
              ports:
                - name: http
                  containerPort: 80
                  protocol: TCP
    ```

    Helm automatically populates the `{{ .Release.Name }}` and `{{ .Values.replicas }}` variables when you render the chart.

1. Run the following command in the root directory to render the chart:

    ```
    helm template --name example .
    ```

1. Check `stdout` to see if the chart rendered correctly.

### Create a service

A [Kubernetes `Service`](https://kubernetes.io/docs/concepts/services-networking/service/) defines a set of pods and a policy by which to access and communicate with them.

Create a `service.yaml` file in the `templates` directory with the following:

```
apiVersion: v1
kind: Service
metadata:
  name: {{ .Release.Name }}-APP_NAME
  labels:
    app.kubernetes.io/name: APP_NAME
    app.kubernetes.io/instance: {{ .Release.Name }}
spec:
  type: ClusterIP
  ports:
    - port: 80
      targetPort: http
      protocol: TCP
      name: http
  selector:
    app.kubernetes.io/name: APP_NAME
    app.kubernetes.io/instance: {{ .Release.Name }}
```
Helm automatically populates the `{{ .Release.Name }}` variable when you render the chart.

### Set up a VirtualService

A `VirtualService` is a stable endpoint that acts as an internal load balancer to send traffic to your Deployment's pods.

Create a `virtualservice.yaml` file in the `templates` directory with the following:

```
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .Release.Name }}-web
spec:
  hosts:
  - "{{ .Release.Name }}.local.govsandbox.uk"
  gateways:
  - "gsp-gsp-cluster.gsp-system"
  http:
  - route:
    - destination:
        host: {{ .Release.Name }}-web
        port:
          number: 3000

```

## Deploy your app

You can deploy your app once Tech Ops has created your namespace and you have created a:

- Helm chart
- Kubernetes Deployment object
- service
- virtual service

To deploy your app, commit your changes and push them to GitHub.

Check that your app is live at {{ .Release.Name }}.{{ .Values.global.cluster.domain }}.

## View your app in the dashboard
_delete and replace as this does not apply to non-local content_

You can view your app in the [Dashboard](https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/) without needing to go through the Service or Ingress that you set up. You do this by using a proxy to access your Kubernetes cluster.

Run the following to use a proxy to access your Kubernetes cluster:

```
kubectl proxy
```

Kubectl will make your Dashboard available at http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/.

To access your dashboard, see the [GSP documentation on accessing dashboards](https://github.com/alphagov/gsp/blob/master/docs/gds-supported-platform/accessing-dashboard.md).
