# =========================
# 1️⃣ Build stage
# =========================
FROM rust:1.75-alpine AS builder

# Install dependency buat static build
RUN apk add --no-cache musl-dev openssl-dev

WORKDIR /app

# Copy manifest dulu biar layer cache kepake
COPY Cargo.toml Cargo.lock ./

# Dummy build biar dependency ke-cache
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo build --release
RUN rm -rf src

# Copy source code asli
COPY src ./src

# Build final binary
RUN cargo build --release


# =========================
# 2️⃣ Runtime stage
# =========================
FROM alpine:3.19

# Security: bikin user non-root
RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

# Copy binary dari builder
COPY --from=builder /app/target/release/rust_engine /app/rust_engine

# Expose port time-tracker
EXPOSE 9000

# Drop privilege
USER app

# Run service
CMD ["./rust_engine"]
