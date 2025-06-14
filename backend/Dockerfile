# Download image
FROM golang:1.24

# Create /app directory and cd into it
WORKDIR /app

# Copy go.mod and go.sum in /app
COPY go.mod go.sum ./

# Install app dependencies
RUN go mod download

# Copy application code in image
COPY main.go ./

# Build server
RUN CGO_ENABLED=0 GOOS=linux go build -o /main

# This does not actually expose it. It's more of a reminder for devs
EXPOSE 8080

# Entrypoint
CMD ["/main"]
