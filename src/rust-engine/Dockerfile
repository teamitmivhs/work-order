# ===== Build stage =====
FROM rust:1.78-slim AS builder

WORKDIR /app

# Copy manifest dulu biar dependency ke-cache
COPY Cargo.toml ./

# Dummy src biar cargo build jalan
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo build --release
RUN rm -rf src

# Copy source asli
COPY src ./src

# Build binary final
RUN cargo build --release

# ===== Runtime stage =====
FROM debian:bookworm-slim

WORKDIR /app

# Install CA cert (penting kalau nanti HTTP keluar)
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/target/release/rust-engine /usr/local/bin/time-tracker

EXPOSE 9000

CMD ["time-tracker"]
