FROM golang:1.22 AS builder

# Install pandoc for markdown to html conversion
RUN apt-get update && apt-get install -y pandoc

# Set working directory and copy code
WORKDIR /app
COPY . .

# Generate index.html from index.md (in /app)
RUN pandoc -s index.md -o index.html

# Build the Go binary (assuming main.go or go.mod is in /app/go-webring)
WORKDIR /app/go-webring
RUN go build -o go-webring

# Final image — use minimal base image
FROM debian:bookworm-slim

# Install ca-certificates (often needed by Go apps)
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Create a non-root user
RUN useradd -m appuser

# Set working directory
WORKDIR /home/appuser

# Copy binary and generated files from builder stage
COPY --from=builder /app/go-webring/go-webring .
COPY --from=builder /app/index.html .
COPY --from=builder /app/go-webring/list.txt .
COPY --from=builder /app/go-webring/static ./static

# Use non-root user
USER appuser

# Expose port from env variable PORT or default 2857
ENV PORT=2857
EXPOSE $PORT

# Run the binary with flags, bind to 0.0.0.0 and port from env var
CMD ["./go-webring", "-l", "0.0.0.0:${PORT}", "-h", "ring.upraxis.org", "-i", "index.html", "-m", "list.txt", "-v", "validation.log"]
