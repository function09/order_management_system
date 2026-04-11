package orders

import (
	"github.com/function09/order_management_system/server/db"
	"github.com/function09/order_management_system/server/internal/addresses"
	"github.com/function09/order_management_system/server/internal/products"
)

type OrderItemInput struct {
	ProductID int
	Quantity  int
}

type AddressInput struct {
	StreetLine1 string
	StreetLine2 string
	City        string
	State       string
	ZipCode     string
}

type SalesOrderInput struct {
	CustomerID  int
	Fulfillment string
	OrderItems  []*OrderItemInput
	Address     *AddressInput
}

type Service struct {
	orderStore    OrderStore
	productsStore products.ProductStore
	addressStore  addresses.AddressStore
	transactor    db.Transactor
}
