# Stage 1: Build the Go application
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go.mod first to define the module
COPY go.mod ./

# Copy all source code, including subdirectories
COPY . .

# Build the application for both architectures
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -o errorpage

# Stage 2: Create the runtime image
FROM alpine:3.20

WORKDIR /app

# Copy the built binary and UI files from the builder stage
COPY --from=builder /app .

# Ensure the binary is executable
RUN chmod +x errorpage

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

EXPOSE 8080

# Run the application
CMD ["./errorpage"]