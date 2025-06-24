# --- Build stage ---
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Install security updates
RUN apk update && apk upgrade && apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o fiscaflow ./cmd/server/main.go

# --- Runtime stage ---
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app

# Copy binary and necessary files
COPY --from=builder /app/fiscaflow /app/fiscaflow
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/.env.example /app/.env.example
COPY --from=builder /app/scripts /app/scripts
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Use non-root user
USER 65532:65532

# Expose app port
EXPOSE 8080

# Healthcheck
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s CMD ["/app/fiscaflow", "-health-check"] || exit 1

# Start application
ENTRYPOINT ["/app/fiscaflow"] 