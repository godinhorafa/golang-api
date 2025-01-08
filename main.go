package main

import (
	"log"
	"net/http"
	"produto-api/internal/database"
	"produto-api/internal/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Carregar variáveis de ambiente do arquivo .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	// Conectar ao banco de dados
	if err := database.Connect(); err != nil {
		log.Fatalf("Erro ao conectar no banco de dados: %v", err)
	}
	defer database.DB.Close()

	// Configurar rotas da API
	r := mux.NewRouter()
	handlers.SetupRoutes(r)

	// Iniciar o servidor
	log.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
