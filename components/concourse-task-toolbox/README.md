
# govsvc/task-toolbox

## Overview

Image containing tools for concourse task scripts.

## Versioning

Please bump `VERSION` file

## Building

```
docker build -t govsvc/task-toolbox:$(cat VERSION) .
```

## Releasing

```
docker push govsvc/task-toolbox:$(cat VERSION)
```
