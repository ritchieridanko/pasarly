# ---------- Build Stage ----------
FROM golang:1.24.2-alpine3.20 AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set work directory
WORKDIR /app/services/notification

# Copy and download app dependencies
COPY shared ../../shared
COPY services/notification/go.mod services/notification/go.sum ./
RUN go mod download

# Copy app source
COPY services/notification/cmd/app ./cmd/app
COPY services/notification/configs ./configs
COPY services/notification/internal ./internal

# Build app
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/app cmd/app/main.go

# ---------- Runtime Stage ----------
FROM alpine:3.20

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Set work directory
WORKDIR /root

# Copy from the Build Stage
COPY --from=builder /app/services/notification/bin ./bin
COPY --from=builder /app/services/notification/configs ./configs

# Set entry point
ENTRYPOINT ["./bin/app"]
