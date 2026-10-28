FROM golang:1.27.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -o tailscale_fwdr -ldflags="-w -s" ./cmd/tailscale-forwarder

FROM busybox:musl

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /app/tailscale_fwdr /usr/local/bin/tailscale_fwdr

ENTRYPOINT ["/usr/local/bin/tailscale_fwdr"]
