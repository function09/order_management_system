package addresses

import (
	"context"
	"database/sql"
)

type Address struct {
	ID          int
	StreetLine1 string
	StreetLine2 string
	City        string
	State       string
	ZipCode     string
	AddressType string
	IsDefault   bool
	CustomerID  int
}

type Store struct {
	*sql.DB
}

type AddressStore interface {
	GetCustomerAddresses(ctx context.Context, cid int) ([]*Address, error)
	GetCustomerAddress(ctx context.Context, aid int) (*Address, error)
	AddCustomerAddress(ctx context.Context, address *Address) (int, error)
	RemoveCustomerAddress(ctx context.Context, aid int) error
}

func (s *Store) GetCustomerAddresses(ctx context.Context, cid int) ([]*Address, error) {
	rows, err := s.QueryContext(ctx, "SELECT street_line_1, street_line_2, city, state, zip_code, address_type, is_default,customer_id from addresses WHERE customer_id=$1", cid)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var addresses []*Address

	for rows.Next() {
		var address Address

		if err := rows.Scan(&address.StreetLine1, &address.StreetLine2, &address.City, &address.State, &address.ZipCode, &address.AddressType, &address.IsDefault, &address.CustomerID); err != nil {
			return nil, err
		}
		addresses = append(addresses, &address)
	}
	return addresses, nil
}
