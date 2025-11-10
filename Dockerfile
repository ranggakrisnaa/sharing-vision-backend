# Multi-stage build for Go Fiber API
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Install build deps
RUN apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Cache modules first
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate cmd/migrate/main.go

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates

# Copy binary
COPY --from=builder /app/server /app/server
COPY --from=builder /app/migrate /app/migrate
COPY --from=builder /app/migrations /app/migrations

ENV PORT=8080
EXPOSE 8080

CMD ["/app/server"]