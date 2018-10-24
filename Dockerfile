FROM golang:alpine AS builder
ADD . /build
ENV CGO_ENABLED=0
WORKDIR /build
RUN go version && go build -mod=vendor


FROM alpine

COPY --from=builder /build/keyguard /app/
COPY loader.sh /app/

RUN apk update && apk add openssh-keygen ca-certificates

ENV PORT 8000
EXPOSE 8000

WORKDIR /app
CMD /app/keyguard
