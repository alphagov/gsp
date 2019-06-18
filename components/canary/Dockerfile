FROM golang as builder

ADD . /go/src/github.com/alphagov/gsp-canary
WORKDIR /go/src/github.com/alphagov/gsp-canary

RUN go get ./... && \
    CGO_ENABLED=0 GOOS=linux go build -o canary -ldflags "-X main.BuildTimestamp=`date +%s`" .

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /go/src/github.com/alphagov/gsp-canary/canary /app/
WORKDIR /app
EXPOSE 8081
CMD ["./canary"]
