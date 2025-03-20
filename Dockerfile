# Usar uma imagem oficial do Go como base
FROM golang:1.20 AS builder

# Definir o diretório de trabalho dentro do container
WORKDIR /app

# Copiar os arquivos do projeto para dentro do container
COPY . .

# Baixar as dependências e compilar o código
RUN go mod download && go build -o bot-financeiro ./cmd/main.go

# Criar uma imagem final menor para rodar o bot
FROM debian:bullseye-slim

WORKDIR /root/

# Copiar o binário compilado da etapa anterior
COPY --from=builder /app/bot-financeiro .

# Definir a variável de ambiente (se necessário)
ENV TELEGRAM_BOT_TOKEN="7984516597:AAFGo8Fceb2mUVyYLsb733W8-fnZozMRRqk"

# Expor portas (se necessário, mas para um bot não costuma ser obrigatório)
# EXPOSE 8080

# Comando para executar o bot
CMD ["./bot-financeiro"]
