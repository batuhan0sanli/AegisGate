# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/aegisgate ./cmd/aegisgate

# Final stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN adduser -D -H -h /app aegisgate

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/aegisgate /app/aegisgate

# Copy entrypoint script
COPY build/Docker/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Create config directory
RUN mkdir -p /app/config && \
    chown -R aegisgate:aegisgate /app

# Switch to non-root user
USER aegisgate

# Expose default port
EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["/app/entrypoint.sh"] 