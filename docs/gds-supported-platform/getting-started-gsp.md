# Getting started with a GDS Supported Platform cluster

## Before you start

You should have:

- a [Docker image](https://docs.docker.com/engine/reference/commandline/images/) of your app built in line with the [12 factor principles](https://docs.cloud.service.gov.uk/architecture.html#12-factor-application-principles)
- access to a [Kubernetes cluster](https://github.com/alphagov/gsp-terraform-ignition/blob/master/docs/gds-supported-platform/troubleshooting_app_errors.md) created by Tech Ops

## Create a GitHub repository

You must create a repository on GitHub to store your Deployment configuration. The GDS Supported Platform only works with repositories stored on GitHub.

## Request namespace

Your service team's Site Reliability Engineer (SRE) must [ask GDS Tech Ops](re-GSP-team@digital.cabinet-office.gov.uk) for a new namespace. You must provide the following information:

* your [GitHub repository address](https://help.github.com/en/articles/about-remote-repositories)
* the name of the [namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/) you want to create
* the email addresses, GitHub usernames, and [GPG public keys](https://www.gnupg.org/gph/en/manual/c14.html) (if applicable) of the users who should get read-only access to the namespace

## Create a Helm chart

The GDS Supported Platform uses a packaging format called [Helm charts](https://helm.sh/docs/developing_charts/). A chart is a collection of files that describe a related set of Kubernetes resources.

You create Helm charts as files in a directory. These files are then packaged into versioned archives that users can deploy.

1. Create a root directory in your GitHub repository. This directory will contain the chart.

1. Create a `Chart.yaml` file in the root directory with the following code:

    ```
    apiVersion: v1
    appVersion: "1.0"
    description: CHART_DESCRIPTION
    name: CHART_NAME
    version: 0.1.0
    ```

    This file defines metadata about the chart.

1. Create a `templates` directory in the root directory. This directory contains all Kubernetes object definitions.

1. Create a `values.yaml` file in the root directory. This file sets the default values for your desired chart variables.

## Create a Kubernetes Deployment object

You run an app by creating a [Kubernetes Deployment object](https://kubernetes.io/docs/concepts/#kubernetes-objects). This object defines your app and its routes, databases and all other relevant information. You describe a Deployment in a YAML file.

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

## Create a service

By default, your apps are not accessible to the public. To expose them to the public, you must set up a [Service](https://kubernetes.io/docs/concepts/services-networking/service/) and an [Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/) into the Kubernetes cluster.

Setting up a Service creates a stable endpoint to send public internet traffic to. To set up a service, create a `service.yaml` file in the `templates` directory with the following code:

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

## Create an Ingress

You must define an [Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/) to route public internet traffic to the stable endpoint you created when you set up the Service.

Create a `ingress.yaml` file in the `templates` directory with the following code:

```
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: {{ .Release.Name }}-myapp
  annotations:
      nginx.ingress.kubernetes.io/rewrite-target: "/"
spec:
  rules:
  - host:  {{ .Release.Name }}.{{ .Values.global.cluster.domain }}
    http:
      paths:
      - backend:
          serviceName: example-myapp
          servicePort: 80
        path: /
```

## Deploy your app

You can deploy your app once Tech Ops has created your namespace and you have created a:

- Helm chart
- Kubernetes Deployment object
- Service
- Ingress

To deploy your app, commit your changes and push them to GitHub.

Check that your app is live at {{ .Release.Name }}.{{ .Values.global.cluster.domain }}.

## View your app in the dashboard

You can view your app in the [Dashboard](https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/) without needing to go through the Service or Ingress that you set up. You do this by using a proxy to access your Kubernetes cluster.

Run the following to use a proxy to access your Kubernetes cluster:

```
kubectl proxy
```

Kubectl will make your Dashboard available at http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/.

To access your dashboard, see the [GSP documentation on accessing dashboards](https://github.com/alphagov/gsp-terraform-ignition/blob/master/docs/gds-supported-platform/accessing-dashboard.md).
