# Etapa de construção do binário
FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bot-financeiro ./cmd/main.go

# Etapa final: executável em uma imagem menor
FROM debian:bookworm-slim

WORKDIR /root/

COPY --from=builder /app/bot-financeiro .

CMD ["./bot-financeiro"]
