package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // Importa o driver do PostgreSQL
)

var DB *sql.DB

func ConectarDB() *sql.DB {
	if DB != nil {
		return DB
	}

	// Pegando as variáveis de ambiente
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSLMODE")

	// Criando a string de conexão
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbName, sslMode,
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Banco inacessível: %v", err)
	}

	fmt.Println("Conectado ao banco com sucesso!")
	return DB
}
