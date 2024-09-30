# Use the official Golang image as a build stage
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /aichatbot

# Copy go.mod and go.sum files first to leverage Docker cache
COPY go.mod go.sum ./

# Download all dependencies (to be cached)
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o myapp .

# Start a new stage from scratch
FROM alpine:latest

# Install necessary dependencies for running the Go app
RUN apk add --no-cache ca-certificates

# Set the Current Working Directory
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /aichatbot/myapp .

# Expose port if needed (e.g., 8080)
EXPOSE 8080

# Command to run the executable
CMD ["./myapp"]
