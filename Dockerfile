# Stage 1: Build
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy .env file if exists (optional)
# COPY .env .

# Copy templates folder
COPY --from=builder /app/templates ./templates

# Expose port
EXPOSE 8080


# Run the application
CMD ["./main"]
