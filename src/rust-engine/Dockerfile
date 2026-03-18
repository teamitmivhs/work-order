# ===== Build stage =====
FROM rust:1.78-slim AS builder

WORKDIR /app

# Copy manifest DAN Cargo.lock dulu biar dependency ke-cache dengan versi yang terkunci
# FIX: tanpa Cargo.lock, setiap build bisa pakai versi dependency berbeda (tidak reproducible)
COPY Cargo.toml Cargo.lock ./

# Dummy src agar cargo bisa download & compile semua dependency dulu
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo build --release
RUN rm -rf src

# Copy source asli lalu build final
COPY src ./src

# FIX: touch main.rs agar cargo tahu file berubah dan tidak skip recompile
RUN touch src/main.rs && cargo build --release

# ===== Runtime stage =====
FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/target/release/rust-engine /usr/local/bin/time-tracker

EXPOSE 9000

# FIX: healthcheck menggunakan /health endpoint yang sudah ada
HEALTHCHECK --interval=15s --timeout=5s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:${PORT:-9000}/health || exit 1

CMD ["time-tracker"]