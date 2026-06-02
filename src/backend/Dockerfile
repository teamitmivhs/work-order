# ===== Build stage =====
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Download dependencies dulu (cache layer terpisah)
COPY go.mod go.sum ./
RUN go mod download

# Copy source dan build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /work-order-backend .

# ===== Runtime stage =====
FROM alpine:latest

WORKDIR /app

# wget untuk healthcheck, ca-certificates untuk HTTPS
RUN apk add --no-cache ca-certificates wget

COPY --from=builder /work-order-backend /work-order-backend

EXPOSE 8080

# Tidak perlu wait-for-db.sh — docker compose sudah handle
# dependency ordering via healthcheck condition di docker-compose.yml
CMD ["/work-order-backend"]