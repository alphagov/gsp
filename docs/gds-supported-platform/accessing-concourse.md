# Accessing Concourse

Each CI cluster that comes with a `ci-system` chart pre installed that will
contain Concourse setup.

For convenience and security, the ingress is disabled by default. This means
Concourse is not accessible via the public internet connection.

You can setup a connection tunnel between your machine and the kubernetes
cluster. This will allow you to access the Concourse Web UI as well as the CLI
endpoint with `fly`.

## Setting up the tunnel

```sh
kubectl port-forward service/gsp-concourse-web -n gsp-system 8080:8080
```

You should be able to visit [`http://127.0.0.1:8080/`](http://127.0.0.1:8080/)
and see the user interface.

The username and password by default are `admin:password`. Based on the need to
setup the tunnel, we believe these do not need to be changed.

## Using `fly`

We strongly discourage the use of `fly`. The Build and Run team are working
hard to provide a consistent way of auto-deploying and updating the pipelines
with a GitOps solution.

The above means, any manual change to the CI system with `fly`, will be
overwritten by an operator.

You may find a need of going against the recommendation, for instance you'd like
to `hijack` a container to debug a failed `job`.

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
          localUser: admin
```

## Accessing for Testing (e.g. Verify concourse pipeline)
It is critical that your local version of 'fly' matches the destination.

Download fly (cli apple icon)
From concourse website:

https://cd.gds-reliability.engineering/teams/gsp/pipelines/verify

```
sudo mv fly /usr/local/bin
sudo chmod +x /usr/local/bin/fly
```

Login

`note:` token will be downloaded automatically, if not, you've got something wrong. You should get a browser uri with a redirect link in the query string. You should be redirected back to the command line session and see `target saved`.

fly -t gsp login -c https://cd.gds-reliability.engineering/ -n gsp

`reponse from site:`

```
logging in to team 'gsp'

navigate to the following URL in your browser (cmd+click on Mac):

https://cd.gds-reliability.engineering/sky/login?redirect_uri=http://127.0.0.1:57517/auth/callback

```
Create a github token with no permissions.

Go to github.com/<your username> -> Settings -> Developer settings -> Personal access tokens -> Create new token with no privileges. Let's call the value "my-secret-value". This will be used for the concourse variable: github-api-token

Get the value of --key-id in the following from the info pipeline (click on the button in https://cd.gds-reliability.engineering/teams/gsp/pipelines/info).

Token created is: "`secret token here`". Note: this is a secret token so must not be publicised, just in case, even though it has no privileges.

```
fly --target gsp hijack --job info/show-available-pipeline-variables sh

aws ssm put-parameter \
   --name "/cd/concourse/pipelines/gsp/my-pipeline-name/my_secret_name" \
   --value "my-secret-value" \
   --type SecureString \
   --key-id "xxx" \
   --overwrite \
   --region eu-west-2
   ```

   This secret will be available in all your pipelines using the syntax: ((github-api-token)).

   Update the pipeline

   ```
  fly -t gsp set-pipeline -p $ACCOUNT_NAME \
    --config tools-staging-prod-infra.yaml \
    --var account-name=$ACCOUNT_NAME \
    --var account-role-arn=$DEPLOYER_ROLE_ARN \
    --var public-gpg-keys=$(yq . ../users/*.yaml | jq -s '[.[] | select(.teams[] | IN("re-gsp")) | .pub]') \
    --check-creds
     ```

Prompts with y/n to accept changes. After accepting, change can be verified with:
     
  fly -t gsp get-pipeline -p verify

See https://docs.aws.amazon.com/cli/latest/reference/ssm/put-parameter.html

