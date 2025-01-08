package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// Connect inicializa a conexão com o banco de dados
func Connect() error {
	// Configuração do DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Abrir conexão com o banco
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("erro ao abrir conexão com o banco: %w", err)
	}

	// Testar a conexão
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("erro ao conectar no banco: %w", err)
	}

	fmt.Println("Conexão com o banco de dados bem-sucedida!")
	return nil
}
