package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"bot-financeiro/internal/bot"
	"bot-financeiro/internal/database"
)

func main() {
	// Apenas carrega o .env se NÃO estiver rodando no Docker
	if os.Getenv("RUNNING_IN_DOCKER") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Aviso: .env não encontrado, usando variáveis do sistema.")
		}
	}

	db := database.ConectarDB()
	defer db.Close()

	database.CriarTabelas(db)

	// Iniciar o bot do Telegram
	bot.IniciarBot()
}
