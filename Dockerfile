FROM golang:alpine as buildbox
WORKDIR /usr/local/app
COPY . .
RUN go build -o go-http-replicator

FROM alpine
COPY --from=buildbox /usr/local/app/go-http-replicator /usr/local/bin/go-http-replicator
