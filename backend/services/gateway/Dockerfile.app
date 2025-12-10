# ---------- Build Stage ----------
FROM golang:1.24.2-alpine3.20 AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set work directory
WORKDIR /app/services/gateway

# Copy and download app dependencies
COPY shared ../../shared
COPY services/gateway/go.mod services/gateway/go.sum ./
RUN go mod download

# Copy app source
COPY services/gateway/internal ./internal
COPY services/gateway/configs ./configs
COPY services/gateway/cmd/app ./cmd/app

# Build app
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/app cmd/app/main.go

# ---------- Runtime Stage ----------
FROM alpine:3.20

# Set work directory
WORKDIR /root

# Copy from the Build Stage
COPY --from=builder /app/services/gateway/bin ./bin
COPY --from=builder /app/services/gateway/configs ./configs

# Expose port
EXPOSE 8080

# Set entry point
ENTRYPOINT ["./bin/app"]
