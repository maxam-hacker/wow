FROM golang:1.22-bookworm AS builder

WORKDIR /build

COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/client

FROM alpine:3.18.2

COPY --from=builder /build/main /

ENTRYPOINT ["/main"]