# syntax=docker/dockerfile:1

FROM golang:bookworm

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code, including local modules
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /userservice

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080

# Run
CMD ["/userservice"]