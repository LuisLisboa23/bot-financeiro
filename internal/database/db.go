package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// ConectarDB inicializa a conexão com o PostgreSQL
func ConectarDB() *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erro ao conectar no banco de dados:", err)
	}

	// Testa a conexão
	err = db.Ping()
	if err != nil {
		log.Fatal("Banco de dados não respondeu:", err)
	}

	fmt.Println("✅ Conexão com PostgreSQL estabelecida!")
	return db
}