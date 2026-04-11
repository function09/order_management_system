package orders

import (
	"context"
	"time"

	"github.com/function09/order_management_system/server/db"
)

type Order struct {
	ID          int
	CustomerID  int
	Status      string
	Fulfillment string
	StreetLine1 string
	StreetLine2 string
	City        string
	State       string
	ZipCode     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OrderItem struct {
	ID        int
	OrderID   int
	ProductID int
	Price     int
	Quantity  int
}

type Store struct {
	dbGetter db.DBGetter
}

type OrderStore interface {
	GetOrders(ctx context.Context, limit int, offset int) ([]*Order, error)
	GetOrder(ctx context.Context, id int) (*Order, error)
	CreateOrder(ctx context.Context, order *Order) (int, error)
	CreateOrderItems(ctx context.Context, orderItems []*OrderItem) error
	UpdateOrderStatus(ctx context.Context, id int, status string) error
}

func (s *Store) GetOrders(ctx context.Context, limit int, offset int) ([]*Order, error) {
	rows, err := s.dbGetter(ctx).QueryContext(ctx, "SELECT id, customer_id, status, fulfillment, street_line_1, street_line_2, city, state, zip_code, created_at, updated_at FROM orders LIMIT $1 OFFSET $2", limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*Order

	for rows.Next() {
		var order Order

		if err := rows.Scan(&order.ID, &order.CustomerID, &order.Status, &order.Fulfillment, &order.StreetLine1, &order.StreetLine2, &order.City, &order.State, &order.ZipCode, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, err
		}

		orders = append(orders, &order)
	}

	return orders, nil
}

func (s *Store) GetOrder(ctx context.Context, id int) (*Order, error) {
	var order Order

	if err := s.dbGetter(ctx).QueryRowContext(ctx, "SELECT id, customer_id, status, fulfillment, street_line_1, street_line_2, city, state, zip_code, created_at, updated_at FROM orders WHERE id=$1", id).Scan(&order.ID, &order.CustomerID, &order.Status, &order.Fulfillment, &order.StreetLine1, &order.StreetLine2, &order.City, &order.State, &order.ZipCode, &order.CreatedAt, &order.UpdatedAt); err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *Store) CreateOrder(ctx context.Context, order *Order) (int, error) {
	var id int
	if err := s.dbGetter(ctx).QueryRowContext(ctx, "INSERT INTO orders (customer_id, status, fulfillment, street_line_1, street_line_2, city, state, zip_code) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id", order.CustomerID, order.Status, order.Fulfillment, order.StreetLine1, order.StreetLine2, order.City, order.State, order.ZipCode).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Store) CreateOrderItems(ctx context.Context, orderItems []*OrderItem) error {
	for _, item := range orderItems {
		_, err := s.dbGetter(ctx).ExecContext(ctx, "INSERT INTO order_items (order_id, product_id, price, quantity) VALUES ($1, $2, $3, $4)", item.OrderID, item.ProductID, item.Price, item.Quantity)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) UpdateOrderStatus(ctx context.Context, id int, status string) error {
	_, err := s.dbGetter(ctx).ExecContext(ctx, "UPDATE orders SET status=$1 WHERE id=$2", status, id)

	return err
}
