# Build the manager binary
FROM golang:1.12.5 as builder

WORKDIR /workspace


ENV KUBEBUILDER_VERSION=2.0.0
RUN wget -nv https://github.com/kubernetes-sigs/kubebuilder/releases/download/v${KUBEBUILDER_VERSION}/kubebuilder_${KUBEBUILDER_VERSION}_linux_amd64.tar.gz \
	&& tar xvf kubebuilder_${KUBEBUILDER_VERSION}_linux_amd64.tar.gz \
	&& mv kubebuilder_${KUBEBUILDER_VERSION}_linux_amd64 /usr/local/kubebuilder
ENV PATH=$PATH:/usr/local/kubebuilder/bin

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
COPY go.mod go.mod
COPY go.sum go.sum
COPY Makefile Makefile
RUN go mod download \
	&& make controller-gen \
	&& go install github.com/onsi/ginkgo/ginkgo

# Copy the go source
COPY . .

# Build and test
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on
RUN make test
RUN go build -a -o manager main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:latest
WORKDIR /
COPY --from=builder /workspace/manager .
ENTRYPOINT ["/manager"]
