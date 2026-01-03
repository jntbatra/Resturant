# Multi-stage build for optimization
FROM golang:1.25 AS builder

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies with explicit timeout and skip checksum verification
RUN go env -w GOSUMDB=off && \
    timeout 120 go mod download -x || \
    (echo "Download failed, retrying..." && sleep 10 && timeout 120 go mod download -x) || \
    (echo "Using vendor if available..." && true)

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o restaurant-app ./cmd/app/main.go

# Final stage - minimal runtime
FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates postgresql-client && \
    rm -rf /var/lib/apt/lists/*

# Create app directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /build/restaurant-app .

# Copy migrations
COPY migrations ./migrations

# Expose port
EXPOSE 8080

# Run the application
CMD ["./restaurant-app"]

