FROM golang:1.16.3-alpine3.13 AS builder

ADD . /build
WORKDIR /build

RUN apk add --no-cache openssh-keygen && chmod 600 testdata/*

ENV CGO_ENABLED=0
RUN go version && go build -mod=vendor

FROM alpine:latest

RUN apk add --no-cache openssh-keygen && \
    addgroup -g 1000 -S keyguard && \
    adduser -u 1000 -S -s /bin/false -G keyguard keyguard

COPY --from=builder --chown=keyguard:keyguard /build/keyguard /app/
COPY --from=builder /build/loader.sh /app/

WORKDIR /app
USER keyguard

ENTRYPOINT ["/app/keyguard"]
