# Use the official Golang image for building the Go application
FROM golang:1.21.3-alpine AS builder

# Install necessary packages for building the Go application with CGO enabled
RUN apk add --no-cache git gcc musl-dev

# Enable CGO
ENV CGO_ENABLED=1

# Set the Current Working Directory inside the container
WORKDIR /app

# Cache the Go modules dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o app ./cmd/fun-banking

# Use a minimal Alpine Linux image for the final stage
FROM alpine:3.18

# Install the CA certificates and SQLite dependencies
RUN apk add --no-cache ca-certificates sqlite-libs

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/app .
COPY .env .
COPY templates/ ./templates
COPY public ./public

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./app"]
