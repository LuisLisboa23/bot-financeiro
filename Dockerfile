# Etapa 1: Construção do binário
FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bot-financeiro ./cmd/main.go

# Etapa 2: Criar imagem final mais leve
FROM debian:bookworm-slim

WORKDIR /root/

# Instalar certificados SSL para evitar erro de TLS
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/bot-financeiro .

CMD ["./bot-financeiro"]
