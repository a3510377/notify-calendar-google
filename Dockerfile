FROM golang:1.19-buster as builder

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY *.go ./
RUN go build -o start_main

FROM debian:buster-slim

WORKDIR /
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
  ca-certificates && \
  rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/start_main .

CMD ./start_main
