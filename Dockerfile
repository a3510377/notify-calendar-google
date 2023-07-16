FROM golang:1.19-buster as builder

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY *.go ./
RUN go build -v -a -ldflags '-s -w' -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}" -o start_main

FROM alpine
WORKDIR /app
COPY --from=builder /app/start_main .

CMD /app/start_main
