FROM golang:latest AS builder
RUN CGO_ENABLED=0 go get -u -v github.com/lixin9311/tlsrouter

FROM alpine:latest AS dist
RUN apk --no-cache add ca-certificates && mkdir -p /etc/tlsrouter
COPY --from=builder /go/bin/go-tlsrouter /sbin
ENTRYPOINT ["tlsrouter"]