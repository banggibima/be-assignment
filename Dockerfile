# --- Build stage ---
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy all source (including tests)
COPY . .

# Build the binary
RUN go build -o server ./cmd/server/main.go

# --- Final stage ---
FROM alpine:latest AS runtime

WORKDIR /app

# Install PostgreSQL client for migrations, etc.
RUN apk add --no-cache postgresql-client bash

# Copy binary, entrypoint, migrations
COPY --from=builder /app/server .
COPY entrypoint.sh .
COPY migrations ./migrations

# Make entrypoint executable
RUN chmod +x entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]
