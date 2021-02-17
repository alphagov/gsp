FROM golang@sha256:9fdb74150f8d8b07ee4b65a4f00ca007e5ede5481fa06e9fd33710890a624331 as builder

ADD . /go/src/github.com/alphagov/gsp-canary
WORKDIR /go/src/github.com/alphagov/gsp-canary

RUN go get ./... && \
    CGO_ENABLED=0 GOOS=linux go build -o canary -ldflags "-X main.BuildTimestamp=`date +%s`" .

FROM alpine@sha256:08d6ca16c60fe7490c03d10dc339d9fd8ea67c6466dea8d558526b1330a85930
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /go/src/github.com/alphagov/gsp-canary/canary /app/
WORKDIR /app
EXPOSE 8081
CMD ["./canary"]
