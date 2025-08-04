# Use official Go image as base
FROM golang:latest

# Set working directory
WORKDIR /app

# Install pandoc
RUN apt-get update && apt-get install -y pandoc

# Clone the repository
RUN git clone https://git.sr.ht/~jbauer/go-webring .

# Convert markdown to HTML
RUN pandoc -s index.md -o index.html

# Build the Go application
RUN go build

# Expose port (adjust if needed, assuming default port)
EXPOSE 2857

# Run the application
CMD ["./go-webring"]
