package database

import (
	"database/sql"
	"log"
)

// CriarTabelas executa as migrations iniciais
func CriarTabelas(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL,
		amount DECIMAL(10,2) NOT NULL,
		category VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Erro ao criar tabela expenses:", err)
	}

	log.Println("✅ Tabela expenses atualizada com user_id!")
}
