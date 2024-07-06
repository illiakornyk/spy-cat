# Stage 1: Build the Go application
FROM golang:1.22.2-alpine AS builder

# Install build tools
RUN apk add --no-cache git gcc g++ libc-dev

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=1 go build -o spy-cat ./cmd/spy-cat

# Stage 2: Run the Go application
FROM alpine:latest

# Install sqlite3
RUN apk add --no-cache sqlite-libs

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/spy-cat .

# Copy the dev config file
COPY config/dev.yml ./config/dev.yml

# Copy the migrations folder
COPY migrations ./migrations

# Set environment variables
ENV CONFIG_PATH=/root/config/dev.yml
ENV CGO_ENABLED=1

# Expose port 8082 to the outside world
EXPOSE 8082

# Command to run the executable
CMD ["./spy-cat"]
