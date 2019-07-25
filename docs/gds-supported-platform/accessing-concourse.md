# Accessing Concourse

GSP provides [Concourse](http://ci.local.govsandbox.uk) for all your building, testing and deploying needs. You
can interact with it through a browser or with `fly`.

## Using `fly`

You can interact with Concourse using [`fly`](https://concourse-ci.org/fly.html). You must make sure the version of `fly` matches the version of Concourse that is running.

The [documentation for `fly`](https://concourse-ci.org/fly.html) describes how
to login and do things with Concourse.

## Using `fly` and Concourse on GSP local

GSP local provides a Concourse that you can experiment with. After your local
cluster has been created you need to port forward to the Istio ingress in order
to be able to access Concourse:

```
sudo --preserve-env kubectl port-forward service/istio-ingressgateway -n istio-system 80:80
```

You can then access Concourse in a browser at `http://ci.local.govsandbox.uk`
and login with `fly` using (`admin:password`):

```
fly login -t gsp-local -c http://ci.local.govsandbox.uk
```
