FROM golang:1.22-bookworm AS builder

WORKDIR /build

COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/server

FROM alpine:3.18.2

COPY --from=builder /build/main /
COPY --from=builder /build/i,rotot.txt /
COPY --from=builder /build/config/server/config.json /config.json

EXPOSE 9877

ENTRYPOINT ["/main"]