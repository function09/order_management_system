package products

import (
	"context"
	"database/sql"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
)

type FakeStore struct {
	GetAllProductsFn func(ctx context.Context) ([]*Product, error)
	GetProductFn     func(ctx context.Context, id int) (*Product, error)
	AddProductFn     func(ctx context.Context, p *Product) error
	UpdateProductFn  func(ctx context.Context, p *Product) error
	RemoveProductFn  func(ctx context.Context, p *Product) error
}

func (s *FakeStore) GetAllProducts(ctx context.Context) ([]*Product, error) {
	return s.GetAllProductsFn(ctx)
}

func (s *FakeStore) GetProduct(ctx context.Context, id int) (*Product, error) {
	return s.GetProductFn(ctx, id)
}

func (s *FakeStore) AddProduct(ctx context.Context, p *Product) error {
	return s.AddProductFn(ctx, p)
}

func (s *FakeStore) UpdateProduct(ctx context.Context, p *Product) error {
	return s.UpdateProductFn(ctx, p)
}

func (s *FakeStore) RemoveProduct(ctx context.Context, p *Product) error {
	return s.RemoveProductFn(ctx, p)
}

func TestGetAllProducts(t *testing.T) {
	var tests = []struct {
		name  string
		store ProductStore
		want  int
	}{
		{"Returns a list of products", &FakeStore{GetAllProductsFn: func(ctx context.Context) ([]*Product, error) {
			return []*Product{{ID: 1, Name: "Pepsi", Price: 199, Quantity: 2, CategoryID: 1}, {ID: 1, Name: "Coke", Price: 299, Quantity: 2, CategoryID: 1}}, nil
		}}, 200},
		{"Returns an empty list of products", &FakeStore{GetAllProductsFn: func(ctx context.Context) ([]*Product, error) {
			return []*Product{}, nil
		}}, 200},
		{"DB call fails", &FakeStore{GetAllProductsFn: func(ctx context.Context) ([]*Product, error) {
			return nil, errors.New("error db call failed")
		}}, 500},
	}

	for _, e := range tests {
		t.Run(e.name, func(t *testing.T) {

			req := httptest.NewRequest("GET", "/products", nil)
			w := httptest.NewRecorder()

			handler := GetAllProductsHandler(e.store)

			handler(w, req)

			if w.Code != e.want {
				t.Errorf("Got %d want %d", w.Code, e.want)
			}
		})
	}
}

func TestGetProduct(t *testing.T) {
	var tests = []struct {
		name      string
		store     ProductStore
		pathValue string
		want      int
	}{
		{"Returns a single product", &FakeStore{GetProductFn: func(ctx context.Context, id int) (*Product, error) {
			return &Product{ID: 1, Name: "Pepsi", Price: 199, Quantity: 2, CategoryID: 1}, nil
		}}, "1", 200},
		{"Returns no product", &FakeStore{GetProductFn: func(ctx context.Context, id int) (*Product, error) {
			return nil, sql.ErrNoRows
		}}, "2", 404},
	}

	for _, e := range tests {
		t.Run(e.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/product/{id}", nil)
			w := httptest.NewRecorder()

			req.SetPathValue("id", e.pathValue)

			handler := GetProductHandler(e.store)
			handler(w, req)

			if w.Code != e.want {
				t.Errorf("Got %d, want %d", w.Code, e.want)
			}

		})
	}
}

func TestAddProduct(t *testing.T) {
	var tests = []struct {
		name  string
		store ProductStore
		body  string
		want  int
	}{
		{"Successfuly adds a product", &FakeStore{AddProductFn: func(ctx context.Context, p *Product) error {
			return nil
		}}, `{"name": "Pepsi", "price":199,"quantity": 5, "category_id": 1}`, 201},
		{"Does not add malformed json", &FakeStore{AddProductFn: func(ctx context.Context, p *Product) error {
			return nil
		}}, `{"name": "Pepsi", "price":,"quantity": 5, "category_id": 1}`, 400},
		{"DB failure", &FakeStore{AddProductFn: func(ctx context.Context, p *Product) error {
			return errors.New("internal server error")
		}}, `{"name": "Pepsi", "price": 199,"quantity": 5, "category_id": 1}`, 500},
	}

	for _, e := range tests {
		t.Run(e.name, func(t *testing.T) {
			body := strings.NewReader(e.body)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/product", body)

			handler := AddProductHandler(e.store)
			handler(w, req)

			if w.Code != e.want {
				t.Errorf("Got %d, want %d", w.Code, e.want)
			}
		})
	}
}

func TestRemoveProduct(t *testing.T) {
	var tests = []struct {
		name      string
		store     ProductStore
		pathValue string
		want      int
	}{
		{"Successfully remove a product", &FakeStore{RemoveProductFn: func(ctx context.Context, p *Product) error {
			return nil
		}}, "1", 200},
		{"Invalid ID", &FakeStore{RemoveProductFn: func(ctx context.Context, p *Product) error {
			return nil
		}}, "abc", 400},
		{"Product not found", &FakeStore{RemoveProductFn: func(ctx context.Context, p *Product) error {
			return sql.ErrNoRows
		}}, "2", 404},
	}

	for _, e := range tests {
		t.Run(e.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/product/{id}", nil)

			req.SetPathValue("id", e.pathValue)

			handler := RemoveProductHandler(e.store)
			handler(w, req)

			if w.Code != e.want {
				t.Errorf("Got %d, want %d", w.Code, e.want)
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	var tests = []struct {
		name      string
		store     ProductStore
		pathValue string
		body      string
		want      int
	}{
		{"Successfuly adds a product", &FakeStore{UpdateProductFn: func(ctx context.Context, p *Product) error {
			return nil
		}}, "1", `{"name": "Pepsi", "price":199,"quantity": 5, "category_id": 1}`, 200},
		{"Does not add malformed json", &FakeStore{UpdateProductFn: func(ctx context.Context, p *Product) error {
			return nil
		}}, "1", `{"name": "Pepsi", "price":,"quantity": 5, "category_id": 1}`, 400},
		{"Invalid ID", &FakeStore{UpdateProductFn: func(ctx context.Context, p *Product) error {
			return nil
		}}, "abc", `{"name": "Pepsi", "price":199,"quantity": 5, "category_id": 1}`, 400},
		{"Product not found", &FakeStore{UpdateProductFn: func(ctx context.Context, p *Product) error {
			return sql.ErrNoRows
		}}, "2", `{"name": "Pepsi", "price":199,"quantity": 5, "category_id": 1}`, 404},
		{"DB failure", &FakeStore{UpdateProductFn: func(ctx context.Context, p *Product) error {
			return errors.New("internal server error")
		}}, "1", `{"name": "Pepsi", "price": 199,"quantity": 5, "category_id": 1}`, 500},
	}

	for _, e := range tests {
		t.Run(e.name, func(t *testing.T) {
			body := strings.NewReader(e.body)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/product/{id}", body)

			req.SetPathValue("id", e.pathValue)

			handler := UpdateProductHandler(e.store)
			handler(w, req)

			if w.Code != e.want {
				t.Errorf("Got %d, want %d", w.Code, e.want)
			}
		})
	}
}
