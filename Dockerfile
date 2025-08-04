# Use official Go image as base
FROM golang:1.22-alpine

# Install pandoc and other dependencies
RUN apk add --no-cache pandoc git

# Set working directory
WORKDIR /app

# Copy go-webring source code into the container
COPY . .

# Generate the homepage HTML from Markdown (if index.md exists)
RUN if [ -f index.md ]; then pandoc -s index.md -o index.html; fi

# Build the Go binary
RUN go build -o go-webring .

# Expose the default port
EXPOSE 2857

# Set environment variable defaults (override with Coolify env vars)
ENV HOST="ring.upraxis.org"
ENV LISTEN="0.0.0.0:2857"
ENV MEMBERS="list.txt"
ENV INDEX="index.html"
ENV VALIDATIONLOG="validation.log"
ENV CONTACT="contact the admin and let them know what's up"

# Run the application
CMD ["./go-webring", \
     "--host", "$HOST", \
     "--listen", "$LISTEN", \
     "--members", "$MEMBERS", \
     "--index", "$INDEX", \
     "--validationlog", "$VALIDATIONLOG", \
     "--contact", "$CONTACT"]
