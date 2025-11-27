# ---------- Build Stage ----------
FROM golang:1.24.2-alpine3.20 AS builder

WORKDIR /app/services/gateway
RUN apk add --no-cache git
ENV PATH="/go/bin:${PATH}"

COPY shared ../../shared
COPY services/gateway/go.mod ./go.mod
COPY services/gateway/go.sum ./go.sum
RUN go mod download

COPY services/gateway/cmd/app ./cmd/app
COPY services/gateway/configs ./configs
COPY services/gateway/internal ./internal

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/app cmd/app/main.go

# ---------- Runtime Stage ----------
FROM alpine:3.20

RUN apk --no-cache add ca-certificates
WORKDIR /root

COPY --from=builder /app/services/gateway/bin ./bin
COPY --from=builder /app/services/gateway/configs ./configs

EXPOSE 8080
ENTRYPOINT ["./bin/app"]
