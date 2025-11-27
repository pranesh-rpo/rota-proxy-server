# Use Go base image
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod file
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o rota-proxy main.go

# Use minimal base image for final container
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/rota-proxy .

# Expose port
EXPOSE 8000

# Run the binary
CMD ["./rota-proxy"]
