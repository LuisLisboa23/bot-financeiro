package test

import (
	"testing"

	"bot-financeiro/internal/bot"
)

// Testes do bot

func TestIniciarBot(t *testing.T) {
	// Apenas verifica se o bot inicializa sem erros
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Erro ao iniciar o bot: %v", r)
		}
	}()

	go bot.IniciarBot()
}