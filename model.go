package main

import (
	"database/sql"
	"fmt"
	// "log"
)

type product struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Quantity int    `json:"quantity"`
	Price   float64 `json:"price"`
}

func getProducts(db *sql.DB) ([]product, error){
	fmt.Println("endpointhit - getproducts")
	query := "SELECT id, name, quantity, price FROM products"
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	products := []product{}
	for rows.Next() {
		var p product
		err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
		if err != nil {
			return nil, err

			
		}
		products = append(products, p)
	}
	return products, nil
	
}

func (p *product) getProduct(db *sql.DB) error {
    fmt.Println("Endpoint hit: getProduct")

    query := "SELECT name, quantity, price FROM products WHERE id = ?"
    row := db.QueryRow(query, p.ID)

    err := row.Scan(&p.Name, &p.Quantity, &p.Price)
    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("no product found with ID %d", p.ID)
        }
        return fmt.Errorf("error retrieving product: %v", err)
    }

    return nil
}


func (p *product) addProduct(db *sql.DB) error{
	fmt.Println("endpoint hit: addProduct")
	query := fmt.Sprintf("INSERT INTO products (name, quantity, price) VALUES ('%s', %d, %f)", p.Name, p.Quantity, p.Price)

	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err 
	}

	p.ID = int(id)
	return nil



}

func (p *product) updateProduct(db *sql.DB) error {
    query := "UPDATE products SET name = ?, quantity = ?, price = ? WHERE id = ?"
    result, err := db.Exec(query, p.Name, p.Quantity, p.Price, p.ID)
    if err != nil {
        return fmt.Errorf("error updating product: %v", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error fetching rows affected: %v", err)
    }
    if rowsAffected == 0 {
        return fmt.Errorf("no rows updated, product with ID %d not found", p.ID)
    }

    return nil
}

func (p *product) deletProduct(db *sql.DB) error {
	query := "DELETE from products WHERE id = ?"
	_, err := db.Exec(query, p.ID)
	return err
}


