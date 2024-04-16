FROM golang:1.19 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Build the application with debug information
RUN go build -gcflags "all=-N -l" -o slip .

FROM golang:1.19

WORKDIR /app

COPY --from=builder /app/slip .

EXPOSE 8084
# Expose the Delve port
EXPOSE 40000

RUN mkdir -p /var/www/slip/notes && chmod 777 /var/www/slip/notes

# Install Delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Start the application with Delve
CMD ["dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./slip"]