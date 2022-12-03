ARG ALPINE_VERSION=3.10
FROM golang:alpine${ALPINE_VERSION} as buildbox
WORKDIR /usr/local/app
COPY . .
RUN go build -o go-http-replicator

FROM alpine:${ALPINE_VERSION}
COPY --from=buildbox /usr/local/app/go-http-replicator /usr/local/bin/go-http-replicator
