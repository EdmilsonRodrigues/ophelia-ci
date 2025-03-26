# --- Build Stage ---
FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download -x

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o ophelia-ci-server -ldflags '-w -s' ./server/. && \
     GOOS=linux GOARCH=amd64 go build -o ophelia-ci-client -ldflags '-w -s' ./client/.

# --- Final Stage ---
FROM debian:stable-slim
WORKDIR /app

RUN mkdir -p /etc/ophelia-ci/ && mkdir -p /var/lib/ophelia/  && \
     apt-get update && apt-get install -y --no-install-recommends git 

COPY --from=builder /app/ophelia-ci-server /app/ophelia-ci-server
COPY --from=builder /app/ophelia-ci-client /usr/bin/ophelia-ci

ENV OPHELIA_CI_FROM_IMAGE=true

ENTRYPOINT ["/app/ophelia-ci-server"]
