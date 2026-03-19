package customers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

type CustomerInput struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func GetAllCustomersHandler(store CustomerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limitString := r.URL.Query().Get("limit")
		offsetString := r.URL.Query().Get("offset")

		limitInt, err := strconv.Atoi(limitString)

		if err != nil || limitInt <= 0 {
			limitInt = 20
		}

		offsetInt, err := strconv.Atoi(offsetString)

		if err != nil || offsetInt <= 0 {
			offsetInt = 0
		}

		customers, err := store.GetAllCustomers(r.Context(), limitInt, offsetInt)

		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(customers)

	}
}

func GetCustomerHandler(store CustomerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.PathValue("id")
		pathValueInt, err := strconv.Atoi(path)

		if err != nil {
			http.Error(w, "Invalid path value", http.StatusBadRequest)
			return
		}

		customer, err := store.GetCustomer(r.Context(), pathValueInt)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Customer not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(customer)
	}
}

func CreateCustomerHandler(store CustomerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var customer *CustomerInput

		if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
			http.Error(w, "Malformed  JSON", http.StatusBadRequest)
			return
		}

		id, err := store.CreateCustomer(r.Context(), &Customer{FirstName: customer.FirstName, LastName: customer.LastName, Email: customer.Email})

		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(id)
	}
}
