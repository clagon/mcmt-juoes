# Stage 1: Build the Go backend binary
FROM golang:1.24.3-alpine AS builder

WORKDIR /src/backend

# Install necessary build tools for SQLite (cgo)
RUN apk add --no-cache gcc musl-dev

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application with CGO enabled (required for go-sqlite3)
RUN CGO_ENABLED=1 GOOS=linux go build -o server-manager ./main.go

# Stage 2: Create the runtime image with Java JRE
# Using eclipse-temurin for Java 21 JRE, as it's a common modern version for Minecraft
FROM eclipse-temurin:21-jre-alpine

WORKDIR /app/backend

# Copy the compiled Go binary from the builder stage
COPY --from=builder /src/backend/server-manager .

# Ensure necessary directories exist
RUN mkdir -p /app/data /app/server

# Expose the API and WebSocket port
EXPOSE 8080

# Command to run the application
CMD ["./server-manager"]
