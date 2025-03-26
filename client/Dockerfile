# --- Build Stage ---
FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download -x

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ophelia-ci-client -ldflags '-w -s' ./client/.

# --- Final Stage ---
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/ophelia-ci-client /app/ophelia-ci-client

ENV OPHELIA_CI_FROM_IMAGE=true

ENTRYPOINT ["/app/ophelia-ci-client"]
