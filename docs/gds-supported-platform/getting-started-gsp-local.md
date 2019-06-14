# Getting started with a local GDS Supported Platform development cluster

These instructions tell you how to set up a local GDS Supported Platform (GSP) cluster and host an app on that cluster.

This process is resource-intensive and you must set the resource amount used by Docker to the following:
- 4 CPUs
- 8 gigabytes of memory

## Build a local GSP cluster

1. Install [Homebrew](https://brew.sh/)

1. Install dependencies from the Brewfile at the root of the repo:

    ```
    brew bundle
    ```

    After installing `docker-machine-driver-hyperkit` follow instructions on how to grant driver superuser privileges to the hypervisor.

1. Clone the GSP repo:

    ```
    git clone https://github.com/alphagov/gsp.git
    ```

1. Run the following to create a local GSP cluster:

    ```
    ./scripts/gsp-local.sh create
    ```

    If this script is still running after 15 minutes, contact the GSP team using the [re-GSP Slack channel](https://gds.slack.com/messages/CDA7YSP0D).

1. Check that you can access your local GSP cluster using [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/):

    ```
    kubectl get nodes
    ```

    You now have a GSP cluster running locally and can access that cluster using `kubectl`.

## Create a Docker image of your app

1. Clone the GOV.UK Prototype Kit:

    ```
    git clone https://github.com/alphagov/govuk-prototype-kit
    ```

1. Create a `Dockerfile` for the `govuk-prototype-kit` app with the following code:

    ```
    FROM node:8.12-alpine

    ADD . /app
    WORKDIR /app

    ARG COLLECT_USAGE_DATA=true

    RUN npm install
    RUN echo "{\"collectUsageData\": $COLLECT_USAGE_DATA}" > usage-data-config.json

    EXPOSE 3000
    CMD ["npm", "start"]
    ```

1. Build the Docker image:

    ```
    docker build . --tag prototype-kit:latest
    ```

1. Test the Docker image:

    ```
    docker run --publish 3000:3000 prototype-kit:latest
    ```
    This will create a local Docker container with the app running in that container.


1. Build the Docker image in the GSP cluster (gsp-local):

    Talk to cluster 'docker' (not the version you have locally installed)
    ```
    eval $(minikube docker-env -p gsp-local)
    docker ps (should yield screens full of containers)
    docker build . -t prototype-kit:latest
    ```
    You have created a Docker image and copied it into your local GSP cluster.
    

## Create a Helm chart

The GDS Supported Platform uses a packaging format called [Helm charts](https://helm.sh/docs/developing_charts/). A chart is a collection of files that describe a related set of Kubernetes resources.

You create Helm charts as files in a directory. These files are then packaged into versioned archives that users can deploy.

1. Create a `chart/` directory inside `govuk-prototype-kit`.

1. Create a `Chart.yaml` file in this directory with the following code:

    ```
    apiVersion: v1
    appVersion: "1.0"
    description: GOV.UK Prototype Kit
    name: prototype-kit
    version: 0.1.0
    ```

    This file defines metadata about the chart.

1. Create a `templates` directory in the `chart/` directory. This directory contains all Kubernetes object definitions.

### Create a Kubernetes Deployment object

You run an app by creating a [Kubernetes Deployment object](https://kubernetes.io/docs/concepts/#kubernetes-objects). This object defines your app and its routes, databases and all other relevant information. You describe a Deployment in a YAML file.

1. Create a `deployment.yaml` file in the `templates` directory with the following code:

    ```
    apiVersion: apps/v1beta2
    kind: Deployment
    metadata:
      name: {{ .Release.Name }}-web
      labels:
        app.kubernetes.io/name: web
        app.kubernetes.io/instance: {{ .Release.Name }}
    spec:
      replicas: {{ .Values.replicas }}
      selector:
        matchLabels:
          app.kubernetes.io/name: web
          app.kubernetes.io/instance: {{ .Release.Name }}
      template:
        metadata:
          labels:
            app.kubernetes.io/name: web
            app.kubernetes.io/instance: {{ .Release.Name }}
        spec:
          containers:
            - name: prototype-kit
              image: "prototype-kit:latest"
              imagePullPolicy: IfNotPresent
              ports:
                - name: http
                  containerPort: 3000
                  protocol: TCP
    ```

1. Run the following command in the root directory to render the chart and send the output to an `output` directory:

    ```
    mkdir output
    helm template --name prototype-kit --output-dir=output chart
    ```

1. Install the contents of the `output` directory to the GSP cluster:

    ```
    kubectl apply -R -f output/
    ```

1. List the Deployments installed in the GSP cluster:

    ```
    kubectl get deployments
    ```

1. Check that the pods are running:

    ```
    kubectl get pods
    ```

### Create a service

By default, your apps are not accessible to the public. To expose them to the public, you must set up a [`Service`](https://kubernetes.io/docs/concepts/services-networking/service/), [`VirtualService`](https://istio.io/docs/reference/config/networking/v1alpha3/virtual-service/) and `port-forward` to the GSP cluster's ingress [Gateway](https://istio.io/docs/reference/config/networking/v1alpha3/gateway/).

Setting up a `VirtualService` creates a stable endpoint that acts like an internal load balancer to send traffic to your Deployment's pods.

1. Create a `service.yaml` file in the `templates` directory with the following code:

    ```
    apiVersion: v1
    kind: Service
    metadata:
      name: {{ .Release.Name }}-web
      labels:
        app.kubernetes.io/name: web
        app.kubernetes.io/instance: {{ .Release.Name }}
    spec:
      type: ClusterIP
      ports:
        - port: 3000
          targetPort: http
          protocol: TCP
          name: http
      selector:
        app.kubernetes.io/name: web
        app.kubernetes.io/instance: {{ .Release.Name }}
    ```

1. Create a `virtualservice.yaml` file in the `templates` directory with the following code:

    ```
    apiVersion: networking.istio.io/v1alpha3
    kind: VirtualService
    metadata:
      name: prototype-kit
    spec:
      hosts:
      - "prototype-kit.local.govsandbox.uk"
      gateways:
      - "gsp-gsp-cluster.gsp-system"
      http:
      - route:
        - destination:
            host: prototype-kit-web
            port:
              number: 3000

    ```

1. Render your template again:

    ```
    helm template --name prototype-kit --output-dir=output chart
    ```

1. Re-apply the template to the GSP cluster:

    ```
    kubectl apply -R -f output/
    ```

You have successfully created a Helm chart.

## Connect to GOV.UK Prototype Kit

1. Use `kubectl port-forward` to tunnel to the ingress Gateway:

    ```
    sudo --preserve-env kubectl port-forward service/istio-ingressgateway -n istio-system 80:80
    ```

1. Open `http://prototype-kit.local.govsandbox.uk` in your browser.

You have successfully deployed the `govuk-prototype-kit` app in your local GSP cluster.

## Destroy cluster

Run the following command to destroy the cluster:

```
./scripts/gsp-local.sh delete
```
