FROM golang:1.15.2-alpine3.12 AS builder

ADD . /build
WORKDIR /build

ENV CGO_ENABLED=0
RUN go version && go build -mod=vendor


FROM alpine:latest

COPY --from=builder /build/keyguard /app/
COPY --from=builder /build/loader.sh /app/

RUN apk add --no-cache openssh-keygen && \
    addgroup -S keyguard && adduser -S -s /bin/false -G keyguard keyguard

WORKDIR /app
USER keyguard

ENTRYPOINT ["/app/keyguard"]
