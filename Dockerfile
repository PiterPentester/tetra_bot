# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build for ARM64 (Orange Pi 5)
ENV CGO_ENABLED=0 GOOS=linux GOARCH=arm64
RUN go build -ldflags "-s -w" -o tetra ./cmd/tetra

# Final stage
FROM scratch

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy Timezone data (important for local time in reports)
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy binary
COPY --from=builder /app/tetra /tetra

# Expose health check port
EXPOSE 8080

# Run as non-root (1000:1000)
USER 1000:1000

ENTRYPOINT ["/tetra"]
