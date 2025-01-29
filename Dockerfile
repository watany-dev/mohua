FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o mohua

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy the pre-built binary file and configs
COPY --from=builder /app/mohua .
COPY --from=builder /app/configs/pricing.yaml ./configs/

# Install ca-certificates for AWS SDK
RUN apk add --no-cache ca-certificates

# Set entrypoint
ENTRYPOINT ["./mohua"]
