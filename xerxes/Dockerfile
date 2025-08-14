# ---- Stage 1: Builder ----
# Use the official Go image as a builder
FROM golang:1.24.5-alpine AS builder

# Set the working directory
WORKDIR /app

# Install git, ca-certificates, and swag CLI
RUN apk add --no-cache git ca-certificates bash

# Install swag CLI for Swagger generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod files to cache dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Generate Swagger documentation
RUN chmod +x scripts/generate-swagger.sh
RUN ./scripts/generate-swagger.sh

# Build the application for release
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# ---- Stage 2: Runner ----
# Use a minimal alpine image for the runtime
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Copy the generated docs directory for Swagger
COPY --from=builder /app/docs ./docs

# Expose the port the app runs on
EXPOSE 10471

# Set the command to run the application
CMD ["./main"]