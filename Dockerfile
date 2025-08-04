# Stage 1: Build the Go binary
FROM golang:1.22 AS builder

# Set working directory
WORKDIR /app

# Copy Go source code and other necessary files
COPY . .

# Optional: generate index.html from index.md if needed
RUN apt-get update && apt-get install -y pandoc
RUN pandoc -s index.md -o index.html

# Build the go-webring binary
RUN go build -o go-webring

# Stage 2: Create minimal runtime image
FROM debian:bookworm-slim

# Create non-root user for safety
RUN useradd -m appuser

# Set working directory
WORKDIR /home/appuser

# Copy binary and assets from builder
COPY --from=builder /app/go-webring .
COPY --from=builder /app/index.html .
COPY --from=builder /app/list.txt .
COPY --from=builder /app/static ./static

# Use the non-root user
USER appuser

# Expose the app's port
EXPOSE 3000

# Start the server
CMD ["./go-webring", "-l", "0.0.0.0:3000", "-h", "ring.upraxis.org", "-i", "index.html", "-m", "list.txt", "-v", "validation.log"]
