# Stage 1: Build the Go application with CGO enabled
FROM golang:1.23 AS builder

# Enable CGO and set up environment variables for Go
ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code, including resources and output directories
COPY . .

# Install the necessary development libraries for SQLite
RUN apt-get update && apt-get install -y gcc libc6-dev

# Build the application
RUN go build -o telegram-budget-bot ./bot/main.go

# Stage 2: Create a lightweight runtime environment with Debian
FROM debian:bookworm-slim

# Install necessary runtime libraries and CA certificates
RUN apt-get update && apt-get install -y libc6 libsqlite3-0 ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy the built binary and directories from the builder stage
COPY --from=builder /app/telegram-budget-bot /telegram-budget-bot
COPY --from=builder /app/resources /resources
COPY --from=builder /app/output /output

# Expose the port your application listens on (if applicable)
# EXPOSE 8080

# Run the Go binary
ENTRYPOINT ["/telegram-budget-bot"]
