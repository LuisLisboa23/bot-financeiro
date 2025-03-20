package main

import (
	"log"

	"github.com/joho/godotenv"
	"bot-financeiro/internal/bot"
	"bot-financeiro/internal/database"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o .env")
	}

	db := database.ConectarDB()
	defer db.Close()

	database.CriarTabelas(db)

	// Iniciar o bot do Telegram
	bot.IniciarBot()
}
