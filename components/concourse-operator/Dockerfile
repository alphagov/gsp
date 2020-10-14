FROM golang:1.13 as builder

# install dep (required when not vendored)
RUN wget https://github.com/golang/dep/releases/download/v0.5.3/dep-linux-amd64 \
	&& mv dep-linux-amd64 /bin/dep \
	&& chmod +x /bin/dep

# install kubebuilder (required for tests)
RUN wget https://github.com/kubernetes-sigs/kubebuilder/releases/download/v1.0.7/kubebuilder_1.0.7_linux_amd64.tar.gz \
	&& tar xvzf kubebuilder_1.0.7_linux_amd64.tar.gz \
	&& mkdir -p /usr/local \
	&& mv kubebuilder_1.0.7_linux_amd64 /usr/local/kubebuilder
ENV PATH="${PATH}:/usr/local/kubebuilder/bin"

# setup context
WORKDIR /go/src/github.com/alphagov/gsp/components/concourse-operator
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
COPY . .

# install dependencies
RUN sh -c 'if [ -e ./vendor ]; then echo skipping dep ensure as found vendor dir 1>&2; else dep ensure -vendor-only; fi'

# run unit tests
ENV KUBEBUILDER_CONTROLPLANE_START_TIMEOUT=1m
ENV KUBEBUILDER_CONTROLPLANE_STOP_TIMEOUT=1m
RUN go test -v ./pkg/... ./cmd/...

# build manager
RUN go build -a -o manager ./cmd/manager

# CA certs
FROM alpine:3.2 as certs
RUN apk add ca-certificates --update

# Minimal image for controller
FROM alpine:3.2
WORKDIR /root/
COPY --from=builder /go/src/github.com/alphagov/gsp/components/concourse-operator/manager /manager
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/manager"]
