# --- Build stage ---
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o fiscaflow ./cmd/server/main.go

# --- Runtime stage ---
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/fiscaflow /app/fiscaflow
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/.env.example /app/.env.example
COPY --from=builder /app/scripts /app/scripts

# Install migration tool (e.g. golang-migrate)
RUN apk add --no-cache ca-certificates curl bash
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz -C /usr/local/bin

# Expose app port
EXPOSE 8080

# Healthcheck
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s CMD wget --spider -q http://localhost:8080/health || exit 1

# Entrypoint (run migrations, then start app)
ENTRYPOINT ["/bin/sh", "-c", "migrate -path /app/migrations -database \"$$DATABASE_URL\" up && exec /app/fiscaflow"] 