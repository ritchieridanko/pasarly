# ---------- Build Stage ----------
FROM golang:1.24.2-alpine3.20 AS builder

WORKDIR /app/services/auth
RUN apk add --no-cache git
ENV PATH="/go/bin:${PATH}"

COPY shared ../../shared
COPY services/auth/go.mod ./go.mod
COPY services/auth/go.sum ./go.sum
RUN go mod download

COPY services/auth/cmd/app ./cmd/app
COPY services/auth/configs ./configs
COPY services/auth/internal ./internal

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/app cmd/app/main.go

# ---------- Runtime Stage ----------
FROM alpine:3.20

RUN apk --no-cache add ca-certificates
WORKDIR /root

COPY --from=builder /app/services/auth/bin ./bin
COPY --from=builder /app/services/auth/configs ./configs

EXPOSE 50051
ENTRYPOINT ["./bin/app"]
