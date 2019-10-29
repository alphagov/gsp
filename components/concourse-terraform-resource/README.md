
# govsvc/terraform-resource

## Overview

Concourse resource for running `terraform apply`.

Extends the [upstream terraform resource](https://github.com/ljfranklin/terraform-resource) to include `awscli`, `git` and `zip` binaries required by `local-exec` scripts.

## Versioning

Please bump `VERSION` file

## Building

```
docker build -t govsvc/terraform-resource:$(cat VERSION) .
```

## Releasing

```
docker push govsvc/terraform-resource:$(cat VERSION)
```
