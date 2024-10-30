# Use the official Go image for building the app
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go mod and sum files and download dependencies
COPY go.mod go.sum .env ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app for production
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/app/main.go

# Use a minimal base image for the runtime
FROM alpine:latest

# Set the working directory in the runtime container
WORKDIR /root/

# Set environment variable for production
ENV GO_ENV=production

# Copy the compiled binary from the builder container
COPY --from=builder /app/server .

# Copy config files if needed
COPY config/config.yaml ./config/config.yaml

# Expose the port your app runs on
EXPOSE 8080

# Run the binary
CMD ["./server"]
