# Build stage
FROM golang:bookworm AS builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY src/go.mod src/go.sum ./
RUN go mod download

# Copy the entire source code, including local modules
COPY src/ .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /userservice

# Final stage
FROM alpine:latest

# Set destination for COPY
WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /userservice userservice 
COPY --from=builder /app/migrations migrations

# Optional: Document the port the application will listen on
EXPOSE 8082

# Run
CMD ["/app/userservice"]