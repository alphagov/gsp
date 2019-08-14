# Set up a local GDS Supported Platform instance

The GDS Supported Platform (GSP) is a platform to host Government Digital Service (GDS) services. The GSP is based on [Docker](https://docs.docker.com/) and [Kubernetes](https://kubernetes.io/).

These instructions explain how to set up a local instance of the GSP and run the example `govuk-prototype-kit` app on that instance. You can use this process to test your app before deploying that app to a production environment on the GSP. This process has several stages. 

1. Install [prerequisite software](#install-prerequisite-software).
1. [Create and run](#create-and-run-a-local-gsp-instance) a local GSP instance.
1. Create a [Docker image of the `govuk-prototype-kit` app](#create-a-docker-image-of-the-govuk-prototype-kit-app).
1. Create a [Helm chart](#create-a-helm-chart).
1. [Deploy the app](#deploy-the-app-on-to-the-local-GSP-instance) on to the local GSP instance.
1. [Connect](#connect-to-the-app) to the app.
1. [Destroy](#destroy-the-local-instance) the local instance.

Setting up a local GSP instance should take no more than 2 hours. Contact the GSP team using the [#re-gsp Slack channel](https://gds.slack.com/messages/CDA7YSP0D) if this process takes longer.

<%= warning_text('While setting up or running a local GSP instance, do not:<br>- power down your laptop<br>- put your laptop in standby mode<br>- connect or disconnect the Cisco VPN client') %>


## Install prerequisite software

1. Install [Homebrew](https://brew.sh/) and run `brew update` to make sure you have the latest version of Homebrew.

1. Run `git clone https://github.com/alphagov/gsp.git` to clone the GSP repository.

1. Go to the local GSP repository folder and run `brew bundle` to install the packages listed in the [Brewfile](https://github.com/alphagov/gsp/blob/master/Brewfile).

1. Grant `driver superuser` privileges to the hypervisor:

    ```
    sudo chown root:wheel /usr/local/bin/docker-machine-driver-hyperkit && sudo chmod u+s /usr/local/bin/docker-machine-driver-hyperkit
    ```

1. Install [Docker](https://docs.docker.com) and configure Docker to have the following settings:

    - CPU = 4
    - Memory = 8 Gb
    - Swap = 2 Gb

1. If you already have either a `~/.minkube/machine/minikube` or `~/.minkube/machine/gocd` directory on your local machine, run the following to make sure Minikube is not running either of these [Kubernetes clusters](https://kubernetes.io/docs/tutorials/kubernetes-basics/create-cluster/cluster-intro/):

    ```
    minikube stop -p CLUSTER
    ```
    where `CLUSTER` can be either `minikube` or `gocd`.

1. If you have already assigned the `KUBECONFIG` environment variable, run `unset KUBECONFIG`.

    
## Create and run a local GSP instance

1. Run `./scripts/gsp-local.sh create` to create a local GSP instance.

    

    This script provisions a [Kubernetes control plane](https://kubernetes.io/docs/concepts/#kubernetes-control-plane) and a curated list of tools on to your local machine.

    This script loops every 10 seconds until the local GSP instance is complete. During the build, you will see looped output similar to the following example:

    ```
    [Stabilize attempt #2] Failed to stabilize. Retrying in 10s...
    kube-system   default-http-backend-5957bfbccb-djj4b       0/1     ErrImagePull       0          104s
    kube-system   metrics-server-6486d4db88-985vz             0/1     ImagePullBackOff   0          103s
    kube-system   nginx-ingress-controller-5bbcd969c5-5wc2c   0/1     ErrImagePull       0          102s
    kube-system   registry-7ngqc                              0/1     ErrImagePull       0          103s
    ```

    Once the process is complete, you will see output similar to the following example:

    ```
    [Apply attempt #1, Stabilize attempt: #3] Finished deploying ./scripts/hack/expose-grafana.yaml.
    Kubernetes master is running at https://192.168.64.9:8443
    KubeDNS is running at https://192.168.64.9:8443/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

    To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
    ```

    If this script is still running after 20 minutes, contact the GSP team using the [#re-gsp Slack channel](https://gds.slack.com/messages/CDA7YSP0D).

1. Check you can access your local GSP instance using [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/):

    ```
    kubectl get nodes
    ```

    You will see output similar to the following example:

    ```
    NAME       STATUS   ROLES    AGE   VERSION
    minikube   Ready    master   14m   v1.12.0
    ```

You now have a local GSP instance you can access using kubectl.




## Create a Docker image of the `govuk-prototype-kit` app

You must build the `govuk-prototype-kit` app as a [Docker image](https://docs.docker.com/v17.09/engine/userguide/storagedriver/imagesandcontainers/) before you can add the app to your GSP local instance.

1. Clone the GOV.UK prototype kit repository:

    ```
    git clone https://github.com/alphagov/govuk-prototype-kit
    ```

1. Go to the root directory of the local `govuk-prototype-kit` repository and create a `Dockerfile` with the following:

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

1. Build a Docker image of the GOV.UK prototype kit app:

    ```
    docker build . --tag prototype-kit:latest
    ```

1. Check that the Dockerised app runs successfully:

    ```
    docker run -it --rm --publish 3000:3000 prototype-kit:latest
    ```
    This will create a local [Docker container](https://www.docker.com/resources/what-container) with the GOV.UK prototype kit app at `http://localhost:3000/`.

    Enter `Ctrl-c` to stop the container once you have confirmed that the GOV.UK prototype kit is running in the local Docker container.

1. Build the Docker image in the GSP instance:

    ```
    eval $(minikube docker-env -p gsp-local)
    docker build . -t prototype-kit:latest
    ```

1. Check that the Docker has finished building the image:

   ```
   docker images | head -n 2
   ```

   Once Docker has built the image, you will see output similar to the following:

    ```
    REPOSITORY      TAG       IMAGE ID        CREATED               SIZE
    prototype-kit   latest    53250c797a93    About a minute ago    246MB
    ```

    You have created a Docker image of the GOV.UK prototype kit app, and built that image in your local GSP instance.


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

### Create the Helm chart directory and Chart.yaml

1. Create a `chart` directory inside the local `govuk-prototype-kit` repository.

1. Create a `Chart.yaml` file in the `chart` directory with the following:

    ```
    apiVersion: v1
    appVersion: "1.0"
    description: GOV.UK Prototype Kit
    name: prototype-kit
    version: 0.1.0
    ```

    This file defines metadata about the Helm chart.

1. Create a `templates` directory in the `chart` directory to store the Kubernetes object definitions.

### Create a Kubernetes Deployment resource

You run an app in the local GSP instance by creating a [Kubernetes Deployment resource object](https://kubernetes.io/docs/concepts/#kubernetes-objects). This object defines your app and its routes, databases and all other relevant information. You describe a Kubernetes Deployment in a `.yaml` file.

1. Create a `deployment.yaml` file in the `templates` directory with the following:

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

1. Run the following command in the root directory of the repository to render the chart and send the output to an `output` directory:

    ```
    mkdir output
    helm template --name prototype-kit --output-dir=output chart
    ```

1. Install the contents of the `output` directory to the local GSP instance:

    ```
    kubectl apply -R -f output/
    ```

    Refer to the [`kubectl apply` documentation](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#apply) for more information.

1. List the Kubernetes Deployments installed in the local GSP instance:

    ```
    kubectl get deployments
    ```

1. [Kubernetes pods](https://kubernetes.io/docs/concepts/workloads/pods/pod/) are the smallest deployable units of computing you can create and manage in Kubernetes. You must check the pods are running:

    ```
    kubectl get pods
    ```

If the pods are running, you have created a Kubernetes Deployment resource.

### Set up a service

A [Kubernetes `Service`](https://kubernetes.io/docs/concepts/services-networking/service/) defines a set of pods and a policy by which to access and communicate with them.

Create a `service.yaml` file in the `templates` directory with the following:

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

### Set up a DestiantionRule

A `DestinationRule` defines policies that apply to traffic intended for a service after routing has occurred.

Create a `destinationrule.yaml` file in the `templates` directory with the following:

```
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: {{ .Release.Name }}-web
  labels:
    app.kubernetes.io/name: web
    app.kubernetes.io/instance: {{ .Release.Name }}
spec:
  host: "{{ .Release.Name }}-web.gsp-system.svc.cluster.local"
  trafficPolicy:
    tls:
      mode: DISABLE
```      

You now have a Helm chart which describes how you will deploy your app to the GSP local instance. 
## Deploy the app on to the local GSP instance

### Create and apply deployable manifests

You must create and apply a set of deployable manifests to tell the GSP local instance how to deploy the `govuk-prototype-kit` app.

1. Go to the root directory of the `govuk-prototype-kit` repository.

1. Render the template:

    ```
    helm template --name prototype-kit --output-dir=output chart
    ```

    You will see output similar to the following example:

    ```
    wrote output/prototype-kit/templates/service.yaml
    wrote output/prototype-kit/templates/deployment.yaml
    wrote output/prototype-kit/templates/destinationrule.yaml
    wrote output/prototype-kit/templates/virtualservice.yaml
    ```

    You now have a set of deployable manifests with the following directory structure and files:

    ```
    output/
    └── prototype-kit/
       └── templates
           ├── deployment.yaml
           ├── destinationrule.yaml
           ├── service.yaml
           └── virtualservice.yaml
    ```

1. Apply the manifests to the local GSP instance:

    ```
    kubectl apply -R -f output/
    ```

    You will see output similar to the following example:

    ```
    deployment.apps/prototype-kit-web created
    destinationrule.networking.istio.io/prototype-kit-web created
    service/prototype-kit-web created
    virtualservice.networking.istio.io/prototype-kit created
    ```

### Check the service, virtualservice, destinatiorule and Kubernetes Deployment are live

1. Check that the Kubernetes Deployment is live:

    ```
    kubectl get deployments prototype-kit-web
    ```

    You will see output similar to the following:

    ```
    NAME                               DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
    prototype-kit-web                  1         0         1            0           1m
    ```
    Once the `CURRENT` and `AVAILABLE` fields are set to `1`, the Kubernetes Deployment is live.

1. Check the `service` is live:

    ```
    kubectl get service prototype-kit-web
    ```

    You will see output similar to the following:

    ```
    NAME                TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
    prototype-kit-web   ClusterIP   10.110.146.100   <none>        3000/TCP   1m
    ```

    Once the `CLUSTER-IP` and `PORT` fields are complete, the service is live.



1. Check the `virtualservice` is live:

    ```
    kubectl get virtualservice prototype-kit-web
    ```
    You will see output similar to the following:

    ```
NAME                GATEWAYS                       HOSTS                                 AGE
prototype-kit-web   [gsp-gsp-cluster.gsp-system]   [prototype-kit.local.govsandbox.uk]   4m
    ```
    Once the `GATEWAYS` and `HOSTS` field have URLs in them, the `virtualservice` is live.

1. Check the `destinationrule` is live:

    ```
    kubectl get destinationrule prototype-kit-web
    ```
    You will see output similar to the following:

    ```
NAME                HOST                                             AGE
prototype-kit-web   prototype-kit-web.gsp-system.svc.cluster.local   8m
    ```
    Once the `HOSTS` field have cluster URL, the `destiantionrule` is live.

Once the `service`, `virtualservice`, `destinationrule` and `Deployment` are live, you have successfully created a local GSP instance and deployed the `govuk-prototype-kit` app.

You must now connect to the app.


## Connect to the app

By default, your app is not accessible outside of the local GSP instance. To connect to the app, you must use [`port-forwarding`](https://kubernetes.io/docs/tasks/access-application-cluster/port-forward-access-application-cluster/).

1. Connect to the `govuk-prototype-kit` app:

    ```
    sudo --preserve-env kubectl port-forward service/istio-ingressgateway -n istio-system 80:80
    ```

    You will see output similar to the following example:

    ```
    sudo --preserve-env kubectl port-forward service/istio-ingressgateway -n istio-system 80:80
    Password:
    Forwarding from 127.0.0.1:80 -> 80
    Forwarding from [::1]:80 -> 80
    ```

1. Open [http://prototype-kit.local.govsandbox.uk](http://prototype-kit.local.govsandbox.uk) in your browser.
1. In the command line, enter `Ctrl-c` to stop port-forwarding.

You have successfully deployed the `govuk-prototype-kit` app in your local GSP instance.


## Destroy the local instance

You should destroy your local GSP instance when you have finished testing your app.

1. Destroy the local instance:

    ```
    ./scripts/gsp-local.sh delete
    ```

1. Remove the installed packages so your system is the same as it was before setting up the local instance:

    ```
    brew uninstall kubernetes-cli
    brew uninstall kubernetes-helm
    brew uninstall hyperkit
    brew uninstall docker-machine-driver-hyperkit
    brew cask uninstall minikube
    ```

You have destroyed your local GSP instance.
