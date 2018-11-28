# Accessing Concourse

Each cluster that comes with a `gsp-base` chart pre installed will contain
Concourse setup.

For convenience and security, the ingress is disabled by default. This means
Concourse is not accessible via the public internet connection.

You can setup a connection tunnel between your machine and the kubernetes
cluster. This will allow you to access the Concourse Web UI as well as the CLI
endpoint with `fly`.

## Setting up the tunnel

```sh
kubectl port-forward service/gsp-base-web -n gsp-base 8080:8080
```

You should be able to visit [`http://127.0.0.1:8080/`](http://127.0.0.1:8080/)
and see the user interface.

The username and password by default are `admin:password`. Based on the need to
setup the tunnel, we believe these do not need to be changed.

## Using `fly`

Concourse on the landing page, should encourage you to download the binary used
to interact with its API.

- [macOS](http://127.0.0.1:8080/api/v1/cli?arch=amd64&platform=darwin)
- [Windows](http://127.0.0.1:8080/api/v1/cli?arch=amd64&platform=windows)
- [GNU/Linux](http://127.0.0.1:8080/api/v1/cli?arch=amd64&platform=linux)

Once downloaded you will be able to login and manage your pipelines.

In order to login run:

```sh
fly -t cluster.staging login -c http://127.0.0.1:8080/
```

The above command will ask you to visit a link in your web browser. This will
authenticate you without providing any login credentials into the CLI.

To push your first pipeline, run:

```sh
fly -t cluster.staging set-pipeline -c ci/build-prototype.yaml -p build-prototype
```

For more information on `fly`, seek the
[documentation](https://concourse-ci.org/fly.html).

## Overcoming default configuration

Sometimes you may wish to allow public access to your Concourse. You can
achieve this by modifying the `values.yaml` file for your cluster and blending
the following into existing setup:

```yaml
concourse:
  concourse:
    web:
      externalUrl: https://concourse.example.com
  web:
    ingress:
      enabled: true
      annotations:
        kubernetes.io/tls-acme: "true"
        kubernetes.io/ingress.class: nginx
      hosts:
      - concourse.example.com
      tls:
      - hosts:
        - concourse.example.com
        secretName: concourse-tls-cert
```

This will establish the URL Concourse would be accessible under, as well as
setup the Let's Encrypt to provide a SSL certificate.

Do remember, the default username and password are easily obtainable. You may
wish to change these too. This can be achieved by adding the following to the
`values.yaml` file:

```yaml
concourse:
  secrets:
    localUsers: admin:password
  concourse:
    web:
      auth:
        mainTeam:
          localUser: "admin"
```
