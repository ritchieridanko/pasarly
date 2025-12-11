# ---------- Build Stage ----------
FROM golang:1.24.2-alpine3.20 AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set work directory
WORKDIR /app/services/user

# Copy and download app dependencies
COPY shared ../../shared
COPY services/user/go.mod services/user/go.sum ./
RUN go mod download

# Copy app source
COPY services/user/cmd/app ./cmd/app
COPY services/user/configs ./configs
COPY services/user/internal ./internal

# Build app
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/app cmd/app/main.go

# ---------- Runtime Stage ----------
FROM alpine:3.20

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Set work directory
WORKDIR /root

# Copy from the Build Stage
COPY --from=builder /app/services/user/bin ./bin
COPY --from=builder /app/services/user/configs ./configs

# Expose port
EXPOSE 50052

# Set entry point
ENTRYPOINT ["./bin/app"]
