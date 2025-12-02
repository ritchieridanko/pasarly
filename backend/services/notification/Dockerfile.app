# ---------- Build Stage ----------
FROM golang:1.24.2-alpine3.20 AS builder

WORKDIR /app/services/notification
RUN apk add --no-cache git
ENV PATH="/go/bin:${PATH}"

COPY shared ../../shared
COPY services/notification/go.mod ./go.mod
COPY services/notification/go.sum ./go.sum
RUN go mod download

COPY services/notification/cmd/app ./cmd/app
COPY services/notification/configs ./configs
COPY services/notification/internal ./internal

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/app cmd/app/main.go

# ---------- Runtime Stage ----------
FROM alpine:3.20

RUN apk --no-cache add ca-certificates
WORKDIR /root

COPY --from=builder /app/services/notification/bin ./bin
COPY --from=builder /app/services/notification/configs ./configs

ENTRYPOINT ["./bin/app"]
