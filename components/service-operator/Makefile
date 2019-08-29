
# Image URL to use all building/pushing image targets
IMG ?= controller:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

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
	ginkgo \
		--failOnPending \
		--failFast \
		--requireSuite \
		--slowSpecThreshold=30 \
		-cover \
		-r \
		./controllers ./apis

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

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
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) webhook paths="./..." output:crd:artifacts:config=config/crd/bases

update-chart: manifests
	cp config/crd/bases/* ../../charts/gsp-cluster/templates/02-gsp-system/service-operator/crd/
	$(CONTROLLER_GEN) rbac:roleName=fixme-up-before-you-gogo paths="./..." output:rbac:stdout | sed 's/fixme-up-before-you-gogo/{{ .Release.Name }}-service-operator-role/g' > ../../charts/gsp-cluster/templates/02-gsp-system/service-operator/role.yaml
	cat config/rbac/role-not-patch.yaml >> ../../charts/gsp-cluster/templates/02-gsp-system/service-operator/role.yaml

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen internal/aws/awsfakes/fake_awsclient.go
	$(CONTROLLER_GEN) object:headerFile=./hack/boilerplate.go.txt paths=./apis/...

internal/aws/awsfakes/fake_awsclient.go:
	go generate ./internal/aws

# Build the docker image
docker-build: test
	docker build . -t ${IMG}
	@echo "updating kustomize image patch file for manager resource"
	sed -i'' -e 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml

# Push the docker image
docker-push:
	docker push ${IMG}

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.0-beta.4
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif
