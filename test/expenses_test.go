package test

import (
	"database/sql"
	"testing"
	"time"

	"bot-financeiro/internal/expenses"

	_ "github.com/lib/pq" // Driver do PostgreSQL
)

func setupTestDB() *sql.DB {
	// Configuração de um banco de dados em memória para testes
	db, err := sql.Open("postgres", "postgres://postgres:230600@localhost:5432/finance_bot?sslmode=disable")
	if err != nil {
		panic(err)
	}

	// Criar tabela para testes
	_, err = db.Exec(`
		CREATE TEMP TABLE expenses (
			id SERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL,
			amount DECIMAL(10,2) NOT NULL,
			category VARCHAR(50) NOT NULL,
			date DATE NOT NULL
		);
	`)
	if err != nil {
		panic(err)
	}

	return db
}

func TestAdicionarDespesa(t *testing.T) {
	db := setupTestDB()
	defer db.Close()

	err := expenses.AdicionarDespesa(db, 1, 100.50, "Transporte", time.Now())
	if err != nil {
		t.Fatalf("Erro ao adicionar despesa: %v", err)
	}
}

func TestListarDespesas(t *testing.T) {
	db := setupTestDB()
	defer db.Close()

	// Adicionar despesas para teste
	_ = expenses.AdicionarDespesa(db, 1, 50.00, "Alimentação", time.Now())
	_ = expenses.AdicionarDespesa(db, 1, 30.00, "Lazer", time.Now())

	despesas, err := expenses.ListarDespesas(db, 1)
	if err != nil {
		t.Fatalf("Erro ao listar despesas: %v", err)
	}

	if len(despesas) != 2 {
		t.Fatalf("Esperado 2 despesas, mas obteve %d", len(despesas))
	}
}

func TestRemoverDespesa(t *testing.T) {
	db := setupTestDB()
	defer db.Close()

	// Adicionar despesa para teste
	_ = expenses.AdicionarDespesa(db, 1, 75.00, "Educação", time.Now())

	// Remover despesa
	err := expenses.RemoverDespesa(db, 1, 1)
	if err != nil {
		t.Fatalf("Erro ao remover despesa: %v", err)
	}

	// Verificar se a despesa foi removida
	despesas, _ := expenses.ListarDespesas(db, 1)
	if len(despesas) != 0 {
		t.Fatalf("Esperado 0 despesas, mas obteve %d", len(despesas))
	}
}
