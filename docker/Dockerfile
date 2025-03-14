FROM golang:1.24-alpine AS builder

LABEL maintainer="marcin_zub@outlook.com"

RUN apk add --no-cache git ca-certificates

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o app \
    ./cmd/main/main.go

# Final stage
FROM scratch

# SSL certificates for HTTPS support (google sheets)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app

COPY --from=builder /build/app .
COPY --from=builder /build/.env.example /app/.env

# non root user
USER 65534:65534

EXPOSE 8080

ENTRYPOINT ["./app"]
