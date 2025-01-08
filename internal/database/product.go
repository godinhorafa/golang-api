package database

import (
	"log"
	"produto-api/internal/models"
	"strconv"
)

type ProductFilters struct {
    Nome      string
    Categoria string
    PrecoMin  string
    PrecoMax  string
    Page      string
    PageSize  string
}

func GetAllProducts() ([]models.Product, error) {
    rows, err := DB.Query("SELECT * FROM produtos")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []models.Product
    for rows.Next() {
        var product models.Product
        if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Category, &product.CreatedAt, &product.UpdatedAt); err != nil {
            return nil, err
        }
        products = append(products, product)
    }

    return products, nil
}

func CreateProduct(product *models.Product) error {
    query := `
        INSERT INTO produtos (name, description, price, category)
        VALUES (?, ?, ?, ?)
    `
    log.Printf("Tentando executar query: %s", query)
    result, err := DB.Exec(query, product.Name, product.Description, product.Price, product.Category)
    if err != nil {
        log.Printf("Erro na execução do SQL: %v", err)
        return err
    }

    id, err := result.LastInsertId()
    if err != nil {
        log.Printf("Erro ao obter último ID: %v", err)
        return err
    }

    product.ID = int(id)
    return nil
}

func UpdateProduct(id string, product *models.Product) error {
    query := `
        UPDATE produtos 
        SET name = ?, description = ?, price = ?, category = ?
        WHERE id = ?
    `
    _, err := DB.Exec(query, product.Name, product.Description, product.Price, product.Category, id)
    if err != nil {
        return err
    }
    return nil
}

func DeleteProduct(id string) error {
    query := "DELETE FROM produtos WHERE id = ?"
    _, err := DB.Exec(query, id)
    if err != nil {
        return err
    }
    return nil
}

func GetProductByID(id string) (*models.Product, error) {
    var product models.Product
    query := "SELECT * FROM produtos WHERE id = ?"
    err := DB.QueryRow(query, id).Scan(
        &product.ID, 
        &product.Name, 
        &product.Description, 
        &product.Price, 
        &product.Category,
        &product.CreatedAt,
        &product.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &product, nil
}

func GetFilteredProducts(filters ProductFilters) ([]models.Product, error) {
    query := "SELECT * FROM produtos WHERE 1=1"
    var args []interface{}

    if filters.Nome != "" {
        query += " AND name LIKE ?"
        args = append(args, "%"+filters.Nome+"%")
    }

    if filters.Categoria != "" {
        query += " AND category = ?"
        args = append(args, filters.Categoria)
    }

    if filters.PrecoMin != "" {
        query += " AND price >= ?"
        args = append(args, filters.PrecoMin)
    }

    if filters.PrecoMax != "" {
        query += " AND price <= ?"
        args = append(args, filters.PrecoMax)
    }

    // Paginação
    page, _ := strconv.Atoi(filters.Page)
    pageSize, _ := strconv.Atoi(filters.PageSize)
    if page <= 0 {
        page = 1
    }
    if pageSize <= 0 {
        pageSize = 10
    }

    offset := (page - 1) * pageSize
    query += " LIMIT ? OFFSET ?"
    args = append(args, pageSize, offset)

    rows, err := DB.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []models.Product
    for rows.Next() {
        var product models.Product
        if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Category, &product.CreatedAt, &product.UpdatedAt); err != nil {
            return nil, err
        }
        products = append(products, product)
    }

    return products, nil
}

func ImportProducts(products []models.Product) error {
    query := `
        INSERT INTO produtos (name, description, price, category)
        VALUES (?, ?, ?, ?)
    `
    
    tx, err := DB.Begin()
    if err != nil {
        return err
    }

    stmt, err := tx.Prepare(query)
    if err != nil {
        tx.Rollback()
        return err
    }
    defer stmt.Close()

    for _, product := range products {
        _, err := stmt.Exec(product.Name, product.Description, product.Price, product.Category)
        if err != nil {
            tx.Rollback()
            return err
        }
    }

    return tx.Commit()
}