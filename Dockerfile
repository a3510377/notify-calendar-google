FROM golang:1.19-buster as builder

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN go build -v -o start_main


FROM debian:buster-slim

WORKDIR /app
COPY --from=builder /app/start_main ./start_main
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
  ca-certificates && \
  rm -rf /var/lib/apt/lists/*

CMD start_main
