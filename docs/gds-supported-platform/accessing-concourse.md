# Accessing Concourse

You use [Concourse](http://ci.local.govsandbox.uk) to build, test and deploy apps on the GDS Supported Platform (GSP). You can access Concourse through either the [`fly` CLI](https://concourse-ci.org/fly.html) or through a browser.

## Accessing Concourse using `fly`

To access Concourse using [`fly`](https://concourse-ci.org/fly.html), your version of `fly` must match your version of Concourse.

Visit your Concourse to check your `fly` version. Use the links in the bottom right of your Concourse page to upgrade your `fly` if necessary. You can find an example Concourse at http://ci.london.verify.govsvc.uk/.

Refer to the [`fly` documentation](https://concourse-ci.org/fly.html) for more information on accessing Concourse using `fly`.

## Accessing Concourse in your local GSP environment using `fly`

The local GSP environment has a Concourse that you can use to try to build, test and deploy your app. 

1. After you [create your local cluster](https://github.com/alphagov/gsp/blob/master/docs/gds-supported-platform/getting-started-gsp-local.md), port forward to the Istio ingress:

    ```
    sudo --preserve-env kubectl port-forward service/istio-ingressgateway -n istio-system 80:80
    ```

1. Access the local environment Concourse at `http://ci.local.govsandbox.uk`:

    ```
    fly login -t gsp-local -c http://ci.local.govsandbox.uk
    ```

1. Sign into `fly` with: 

    - username: admin
    - password: password
