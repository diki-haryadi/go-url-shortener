# Build stage
FROM golang:1.21-alpine AS builder

# Add git for private repos if needed
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

# Final stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose the port
EXPOSE 8080

# Command to run
CMD ["./main"]