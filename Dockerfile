FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -o tailscale_fwdr -ldflags="-w -s" ./.

FROM gcr.io/distroless/static

WORKDIR /app

COPY --from=builder /app/tailscale_fwdr /usr/local/bin/tailscale_fwdr

ENTRYPOINT ["/usr/local/bin/tailscale_fwdr"]