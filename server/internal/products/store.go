package products

import (
	"database/sql"
)

type Product struct {
	ID         int
	Name       string
	Price      int
	Quantity   int
	CategoryID int
}

type Store struct {
	*sql.DB
}

type ProductStore interface {
	GetAllProducts() ([]*Product, error)
	GetProduct(id int) (*Product, error)
	AddProduct(p *Product) error
	RemoveProduct(p *Product) error
	UpdateProduct(p *Product) error
}

func (s *Store) GetAllProducts() ([]*Product, error) {
	rows, err := s.Query("SELECT id, name, price, quantity, category_id FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*Product

	for rows.Next() {
		var prod Product

		if err := rows.Scan(&prod.ID, &prod.Name, &prod.Price, &prod.Quantity, &prod.CategoryID); err != nil {
			return products, err
		}
		products = append(products, &prod)
	}

	return products, nil
}

func (s *Store) GetProduct(id int) (*Product, error) {
	var product Product
	if err := s.QueryRow("SELECT id, name, price, quantity, category_id FROM products WHERE id = $1", id).Scan(&product.ID, &product.Name, &product.Price, &product.Quantity, &product.CategoryID); err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *Store) AddProduct(p *Product) error {
	_, err := s.Exec("INSERT INTO products (name, price, quantity, category_id) VALUES ($1, $2, $3, $4)", p.Name, p.Price, p.Quantity, p.CategoryID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) RemoveProduct(p *Product) error {
	_, err := s.Exec("DELETE FROM products WHERE id = $1", p.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateProduct(p *Product) error {
	_, err := s.Exec("UPDATE products SET name = $1, price = $2, quantity = $3, category_id = $4 WHERE id = $5", p.Name, p.Price, p.Quantity, p.CategoryID)

	if err != nil {
		return err
	}

	return nil
}
