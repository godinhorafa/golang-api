package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"produto-api/internal/database"
	"produto-api/internal/models"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// SetupRoutes configura todas as rotas da API
func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/produtos", GetProducts).Methods("GET")
	r.HandleFunc("/produtos/{id:[0-9]+}", GetProductByID).Methods("GET")
	r.HandleFunc("/produtos", CreateProduct).Methods("POST")
	r.HandleFunc("/produtos/importar", ImportProducts).Methods("POST")
	r.HandleFunc("/produtos/{id:[0-9]+}", UpdateProduct).Methods("PUT")
	r.HandleFunc("/produtos/{id:[0-9]+}", DeleteProduct).Methods("DELETE")
}

// GetProducts retorna todos os produtos
func GetProducts(w http.ResponseWriter, r *http.Request) {
	// Parâmetros de filtro
	filters := database.ProductFilters{
		Nome:      r.URL.Query().Get("nome"),
		Categoria: r.URL.Query().Get("categoria"),
		PrecoMin:  r.URL.Query().Get("preco_min"),
		PrecoMax:  r.URL.Query().Get("preco_max"),
		Page:      r.URL.Query().Get("page"),
		PageSize:  r.URL.Query().Get("page_size"),
	}

	products, err := database.GetFilteredProducts(filters)
	if err != nil {
		http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
    var product models.Product
    if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
        log.Printf("Erro ao decodificar produto: %v", err)
        http.Error(w, "Erro ao decodificar produto", http.StatusBadRequest)
        return
    }

    if err := database.CreateProduct(&product); err != nil {
        log.Printf("Erro ao criar produto: %v", err)
        http.Error(w, "Erro ao criar produto", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(product)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var product models.Product
    if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
        http.Error(w, "Erro ao decodificar produto", http.StatusBadRequest)
        return
    }

    if err := database.UpdateProduct(id, &product); err != nil {
        http.Error(w, "Erro ao atualizar produto", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(product)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    if err := database.DeleteProduct(id); err != nil {
        http.Error(w, "Erro ao deletar produto", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// Novo handler para buscar produto por ID
func GetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	product, err := database.GetProductByID(id)
	if err != nil {
		http.Error(w, "Produto não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// Novo handler para importação em massa
func ImportProducts(w http.ResponseWriter, r *http.Request) {
	// Verificar o Content-Type
	contentType := r.Header.Get("Content-Type")
	var products []models.Product

	if strings.Contains(contentType, "application/json") {
		// Processar JSON
		if err := json.NewDecoder(r.Body).Decode(&products); err != nil {
			http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
			return
		}
	} else if strings.Contains(contentType, "text/csv") {
		// Processar CSV
		reader := csv.NewReader(r.Body)
		records, err := reader.ReadAll()
		if err != nil {
			http.Error(w, "Erro ao ler CSV", http.StatusBadRequest)
			return
		}
		
		// Pular cabeçalho
		for _, record := range records[1:] {
			price, _ := strconv.ParseFloat(record[2], 64)
			product := models.Product{
				Name:        record[0],
				Description: record[1],
				Price:      price,
				Category:   record[3],
			}
			products = append(products, product)
		}
	} else {
		http.Error(w, "Formato não suportado", http.StatusBadRequest)
		return
	}

	// Importar produtos
	if err := database.ImportProducts(products); err != nil {
		http.Error(w, "Erro ao importar produtos", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Importados %d produtos com sucesso", len(products)),
	})
}