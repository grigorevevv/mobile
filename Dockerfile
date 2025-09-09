# Стадия сборки
FROM golang:1.24 AS builder

COPY go.mod go.sum ./
RUN go mod download

WORKDIR /app

COPY . .

ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0

RUN go build -v -o subscription ./cmd/server/main.go

# Финальный образ
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /

COPY --from=builder /app/subscription .

RUN chmod +x /subscription

EXPOSE 8070

CMD ["./subscription"]
