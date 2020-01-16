
# Image URL to use all building/pushing image targets
IMG ?= controller:latest

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: manager update-chart

# Lint
lint: fmt vet
	golangci-lint run

# Run tests
test: generate manifests
	go test -mod=vendor ./controllers ./apis/... ./internal/...

# Build manager binary
manager: generate fmt vet
	go build -mod=vendor -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./main.go

# Install CRDs into a cluster
install: manifests
	kubectl apply -f config/crd/bases

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	kubectl apply -f config/crd/bases
	kustomize build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: generate

update-chart: manifests
	cp config/crd/bases/* ../../charts/gsp-cluster/templates/02-gsp-system/service-operator/crd/
	sed 's/fixme-up-before-you-gogo/{{ .Release.Name }}-service-operator-role/g' < ./config/rbac/fixme/* > ../../charts/gsp-cluster/templates/02-gsp-system/service-operator/role.yaml
	cat config/rbac/role-not-patch.yaml >> ../../charts/gsp-cluster/templates/02-gsp-system/service-operator/role.yaml

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: internal/aws/sdk/sdkfakes/fake_client.go internal/aws/cloudformation/cloudformationfakes/fake_stack.go
	go generate

internal/aws/sdk/sdkfakes/fake_client.go: internal/aws/sdk/client.go
	go generate ./internal/aws/sdk

internal/aws/cloudformation/cloudformationfakes/fake_stack.go: internal/aws/cloudformation/types.go
	go generate ./internal/aws/cloudformation

# Build the docker image
docker-build: test
	docker build . -t ${IMG}
	@echo "updating kustomize image patch file for manager resource"
	sed -i'' -e 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml

# Push the docker image
docker-push:
	docker push ${IMG}

