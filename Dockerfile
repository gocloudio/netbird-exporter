FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o netbird-exporter ./cmd/netbird-exporter

# Run stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/netbird-exporter .

# Expose Prometheus metrics port
EXPOSE 9101

# Run the netbird-exporter binary
ENTRYPOINT ["./netbird-exporter"] 