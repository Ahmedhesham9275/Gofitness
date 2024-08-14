# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o myblog

# Run stage
FROM alpine:3.18
WORKDIR /app

# Copy the executable from the builder stage
COPY --from=builder /app/myblog /app/myblog
COPY --from=builder /app/config /app/config
COPY --from=builder /app/controllers /app/controllers
COPY --from=builder /app/middlewares /app/middlewares
COPY --from=builder /app/models /app/models
COPY --from=builder /app/routes /app/routes
COPY --from=builder /app/utils /app/utils

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["/app/myblog"]
