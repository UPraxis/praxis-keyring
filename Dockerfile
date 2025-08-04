# Stage 1: Build stage
FROM golang:1.22 AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first for caching dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all source files, including index.md, list.txt, static folder, etc.
COPY . .

# Install pandoc for markdown to HTML conversion
RUN apt-get update && apt-get install -y pandoc

# Generate index.html from index.md
RUN pandoc -s index.md -o index.html

# Build the Go binary (adjust main.go if needed)
RUN go build -o go-webring

# Stage 2: Runtime image - small and secure
FROM debian:bookworm-slim

# Create a non-root user for security
RUN useradd -m appuser

# Set working directory for runtime
WORKDIR /home/appuser

# Copy binary and assets from the builder stage
COPY --from=builder /app/go-webring .
COPY --from=builder /app/index.html .
COPY --from=builder /app/list.txt .
COPY --from=builder /app/static ./static

# Change ownership to appuser (optional but recommended)
RUN chown -R appuser:appuser /home/appuser

# Switch to the non-root user
USER appuser

# Expose the port your app listens on
EXPOSE 3000

# Run the binary with the required flags for the ring
CMD ["./go-webring", "-l", "0.0.0.0:3000", "-h", "ring.upraxis.org", "-i", "index.html", "-m", "list.txt", "-v", "validation.log"]
