package customers

import (
	"context"
	"database/sql"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
)

type FakeStore struct {
	GetAllCustomersFn func(ctx context.Context, limit int, offset int) ([]*Customer, error)
	GetCustomerFn     func(ctx context.Context, id int) (*Customer, error)
	CreateCustomerFn  func(ctx context.Context, cst *Customer) (int, error)
	RemoveCustomerFn  func(ctx context.Context, id int) error
	UpdateCustomerFn  func(ctx context.Context, cst *Customer) error
}

func (s *FakeStore) GetAllCustomers(ctx context.Context, limit int, offset int) ([]*Customer, error) {
	return s.GetAllCustomersFn(ctx, limit, offset)
}

func (s *FakeStore) GetCustomer(ctx context.Context, id int) (*Customer, error) {
	return s.GetCustomerFn(ctx, id)
}
func (s *FakeStore) CreateCustomer(ctx context.Context, cst *Customer) (int, error) {
	return s.CreateCustomerFn(ctx, cst)
}

func (s *FakeStore) RemoveCustomer(ctx context.Context, id int) error {
	return s.RemoveCustomerFn(ctx, id)
}
func (s *FakeStore) UpdateCustomer(ctx context.Context, cst *Customer) error {
	return s.UpdateCustomerFn(ctx, cst)
}

func TestGetAllCustomers(t *testing.T) {
	var tests = []struct {
		name  string
		store CustomerStore
		param string
		want  int
	}{
		{"returns a list of customers", &FakeStore{GetAllCustomersFn: func(ctx context.Context, limit, offset int) ([]*Customer, error) {
			return []*Customer{{ID: 1, FirstName: "John", LastName: "Doe", Email: "johndoe@mail.com", IsActive: true}, {ID: 2, FirstName: "Jane", LastName: "Doe", Email: "janedoe@mail.com", IsActive: true}}, nil
		}}, "?limit=5&offset=0", 200},
		{"returns an empty list of customers", &FakeStore{GetAllCustomersFn: func(ctx context.Context, limit, offset int) ([]*Customer, error) {
			return []*Customer{}, nil
		}}, "?limit=5&offset=0", 200},
		{"DB call fails", &FakeStore{GetAllCustomersFn: func(ctx context.Context, limit, offset int) ([]*Customer, error) {
			return nil, errors.New("error db call failed")
		}}, "?limit=5&offset=0", 500},
		{"Invalid limit parameter", &FakeStore{GetAllCustomersFn: func(ctx context.Context, limit int, offset int) ([]*Customer, error) {
			return []*Customer{{ID: 1, FirstName: "John", LastName: "Doe", Email: "johndoe@mail.com", IsActive: true}, {ID: 2, FirstName: "Jane", LastName: "Doe", Email: "janedoe@mail.com", IsActive: true}}, nil
		}}, "?limit=abc&offset=0", 200},
		{"Invalid offset parameter", &FakeStore{GetAllCustomersFn: func(ctx context.Context, limit int, offset int) ([]*Customer, error) {
			return []*Customer{{ID: 1, FirstName: "John", LastName: "Doe", Email: "johndoe@mail.com", IsActive: true}, {ID: 2, FirstName: "Jane", LastName: "Doe", Email: "janedoe@mail.com", IsActive: true}}, nil
		}}, "?limit=5&offset=abc", 200},
		{"No query parameter", &FakeStore{GetAllCustomersFn: func(ctx context.Context, limit int, offset int) ([]*Customer, error) {
			return []*Customer{{ID: 1, FirstName: "John", LastName: "Doe", Email: "johndoe@mail.com", IsActive: true}, {ID: 2, FirstName: "Jane", LastName: "Doe", Email: "janedoe@mail.com", IsActive: true}}, nil
		}}, "?limit=&offset=", 200},
	}
	for _, e := range tests {
		t.Run(e.name, func(t *testing.T) {

			req := httptest.NewRequest("GET", "/customers"+e.param, nil)
			w := httptest.NewRecorder()

			handler := GetAllCustomersHandler(e.store)

			handler(w, req)

			if w.Code != e.want {
				t.Errorf("Got %d want %d", w.Code, e.want)
			}

		})

	}
}

func TestGetCustomer(t *testing.T) {
	var tests = []struct {
		name      string
		store     CustomerStore
		pathValue string
		want      int
	}{
		{"Returns a customer", &FakeStore{GetCustomerFn: func(ctx context.Context, id int) (*Customer, error) {
			return &Customer{ID: 1, FirstName: "John", LastName: "Doe", Email: "johndoe@mail.com", IsActive: true}, nil
		}}, "1", 200},
		{"Returns no customer", &FakeStore{GetCustomerFn: func(ctx context.Context, id int) (*Customer, error) {
			return nil, sql.ErrNoRows
		}}, "2", 404},
		{"No path or bad path", &FakeStore{GetCustomerFn: func(ctx context.Context, id int) (*Customer, error) {
			return nil, errors.New("Bad path")
		}}, "", 400},
		{"Internal server error", &FakeStore{GetCustomerFn: func(ctx context.Context, id int) (*Customer, error) {
			return nil, errors.New("Internal server error")
		}}, "1", 500},
	}
	for _, e := range tests {
		t.Run(e.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/customers/"+e.pathValue, nil)
			w := httptest.NewRecorder()

			req.SetPathValue("id", e.pathValue)
			handler := GetCustomerHandler(e.store)
			handler(w, req)

			if w.Code != e.want {
				t.Errorf("Got %d want %d", w.Code, e.want)
			}
		})
	}
}

func TestCreateCustomer(t *testing.T) {
	var tests = []struct {
		name  string
		store CustomerStore
		body  string
		want  int
	}{
		{"Creates a customer", &FakeStore{CreateCustomerFn: func(ctx context.Context, cst *Customer) (int, error) {
			return 1, nil
		}}, `{"firstName":"John", "lastName":"Doe", "email":"johndoe@mail.com"}`, 201},
		{"Malformed JSON", &FakeStore{CreateCustomerFn: func(ctx context.Context, cst *Customer) (int, error) {
			return 0, nil
		}}, `{"firstName":, "lastName":"Doe", "email":"johndoe@mail.com"}`, 400},
		{"DB Error", &FakeStore{CreateCustomerFn: func(ctx context.Context, cst *Customer) (int, error) {
			return 0, errors.New("Internal server error")
		}}, `{"firstName":"John", "lastName":"Doe", "email":"johndoe@mail.com"}`, 500},
	}

	for _, e := range tests {
		t.Run(e.name, func(t *testing.T) {
			body := strings.NewReader(e.body)
			req := httptest.NewRequest("POST", "/customer", body)
			w := httptest.NewRecorder()

			handler := CreateCustomerHandler(e.store)
			handler(w, req)

			if w.Code != e.want {
				t.Errorf("Got %d want %d", w.Code, e.want)
			}
		})
	}

}
