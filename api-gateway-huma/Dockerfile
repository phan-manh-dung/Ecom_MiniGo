# Build stage
FROM golang:1.24 AS builder

# Set working directory
WORKDIR /app

# Copy go.mod file first
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build only the api-gateway
    RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-gateway-huma .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/api-gateway-huma .

# Expose port
EXPOSE 8081


# Run the binary
CMD ["./api-gateway-huma"]
