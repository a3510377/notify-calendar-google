FROM golang:1.19-buster as builder

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY *.go ./
RUN go build -v -a -ldflags '-s -w' -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}"

FROM alpine
WORKDIR /app
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
  ca-certificates && \
  rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/start_main .

CMD /app/start_main
