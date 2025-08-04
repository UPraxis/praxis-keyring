# Use official Go image as base
FROM golang:latest

# Set working directory
WORKDIR /app

# Install pandoc and git
RUN apt-get update && apt-get install -y pandoc git

# Clone the repository
RUN git clone https://github.com/UPraxis/praxis-keyring.git .

# Optional: Generate HTML from markdown (if index.md exists)
RUN [ -f index.md ] && pandoc -s index.md -o index.html || echo "index.md not found"

# Download Go dependencies
RUN go mod tidy

# Build the Go application with the correct name
RUN go build -o go-webring

# Expose the port (default is 2857 unless overridden)
EXPOSE 2857

# Run the application (split command and arguments)
CMD ["./go-webring", "--host", "ring.upraxis.org"]
