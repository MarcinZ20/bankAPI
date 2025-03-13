FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

RUN addgroup app && adduser -S -G app app

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

# Copy SSL certificates for HTTPS support
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app

COPY --from=builder /build/app .

# non-root alpine user
USER 65532:65532

EXPOSE 8080

ENTRYPOINT ["./app"]
