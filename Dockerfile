FROM golang:alpine AS builder
ADD . /src/github.com/cromega/keyguard
ENV GOPATH /
WORKDIR /src/github.com/cromega/keyguard
RUN go build


FROM alpine

COPY --from=builder /src/github.com/cromega/keyguard/keyguard /app/
COPY loader.sh /app/

RUN apk update && apk add openssh-keygen

ENV PORT 8000
EXPOSE 8000

WORKDIR /app
CMD /app/keyguard
