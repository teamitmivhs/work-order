# Stage 1: Build the Go binary
FROM golang:1.21-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
# -ldflags="-w -s" reduces the size of the binary by stripping debug information
# CGO_ENABLED=0 is important for a static binary that runs in a minimal container like alpine or scratch
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /work-order-backend .

# Stage 2: Create the final, minimal image
FROM alpine:latest

# Install mysql-client for the wait script

RUN apk add --no-cache mysql-client netcat-openbsd

# Set the Current Working Directory inside the container
WORKDIR /

# Copy the pre-built binary from the builder stage
COPY --from=builder /work-order-backend /work-order-backend

# Copy the wait-for-db script
COPY wait-for-db.sh /wait-for-db.sh
RUN chmod +x /wait-for-db.sh

# Expose port 8080 to the outside world
EXPOSE 8080





# Command to run the executable
CMD ["/work-order-backend"]
