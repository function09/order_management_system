package products

import (
	"context"
	"database/sql"
)

type Product struct {
	ID         int
	SKU        string
	Name       string
	Price      int
	Quantity   int
	CategoryID int
}

type Store struct {
	*sql.DB
}

type ProductStore interface {
	GetAllProducts(ctx context.Context, limit, offset int) ([]*Product, error)
	GetProduct(ctx context.Context, id int) (*Product, error)
	AddProduct(ctx context.Context, p *Product) error
	RemoveProduct(ctx context.Context, p *Product) error
	UpdateProduct(ctx context.Context, p *Product) error
}

func (s *Store) GetAllProducts(ctx context.Context, limit int, offset int) ([]*Product, error) {
	rows, err := s.QueryContext(ctx, "SELECT id, sku, name, price, quantity, category_id FROM products ORDER BY id ASC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*Product

	for rows.Next() {
		var prod Product

		if err := rows.Scan(&prod.ID, &prod.SKU, &prod.Name, &prod.Price, &prod.Quantity, &prod.CategoryID); err != nil {
			return products, err
		}
		products = append(products, &prod)
	}

	return products, nil
}

func (s *Store) GetProduct(ctx context.Context, id int) (*Product, error) {
	var product Product
	if err := s.QueryRowContext(ctx, "SELECT id, sku, name, price, quantity, category_id FROM products WHERE id = $1", id).Scan(&product.ID, &product.SKU, &product.Name, &product.Price, &product.Quantity, &product.CategoryID); err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *Store) AddProduct(ctx context.Context, p *Product) error {
	_, err := s.ExecContext(ctx, "INSERT INTO products (name, price, quantity, category_id, sku) VALUES ($1, $2, $3, $4, $5)", p.Name, p.Price, p.Quantity, p.CategoryID, p.SKU)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) RemoveProduct(ctx context.Context, p *Product) error {
	_, err := s.ExecContext(ctx, "DELETE FROM products WHERE id = $1", p.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateProduct(ctx context.Context, p *Product) error {
	_, err := s.ExecContext(ctx, "UPDATE products SET name = $1, price = $2, quantity = $3, category_id = $4, sku = $5 WHERE id = $6", p.Name, p.Price, p.Quantity, p.CategoryID, p.SKU, p.ID)

	if err != nil {
		return err
	}

	return nil
}
