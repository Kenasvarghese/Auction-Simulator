# ---- Build stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the binary from the main entrypoint
RUN go build -o auction-simulator ./app/main.go

# ---- Runtime stage ----
FROM alpine:latest

WORKDIR /app

# Copy binary only (not the source)
COPY --from=builder /build/auction-simulator .


# Default command
CMD ["./auction-simulator"]
