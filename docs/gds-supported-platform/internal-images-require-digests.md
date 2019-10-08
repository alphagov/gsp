# Internal images 

We enforce using an immutable digest for the image in order to avoid the situation where a Docker image can be replaced without being built by Concourse (for example, by pushing to the registry with some mutable tag which is referenced by a Helm chart).

This means you must reference images from the internal registry like this:

```
registry.local.govsandbox.uk/image@sha256:01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b
```

Referencing images in the internal registry will fail if a SHA256 is not used:

```
# This will fail
registry.local.govsandbox.uk/image:latest
```
