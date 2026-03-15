package products

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

type ProductInput struct {
	Name       string `json:"name"`
	SKU        string `json:"sku"`
	Price      int    `json:"price"`
	Quantity   int    `json:"quantity"`
	CategoryID int    `json:"category_id" `
}

func GetAllProductsHandler(store ProductStore) http.HandlerFunc {
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

		products, err := store.GetAllProducts(r.Context(), limitInt, offsetInt)

		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(products)
	}
}

func GetProductHandler(store ProductStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathValue := r.PathValue("id")
		pathValueInt, err := strconv.Atoi(pathValue)

		if err != nil {
			http.Error(w, "Invalid path value", http.StatusBadRequest)
			return
		}

		product, err := store.GetProduct(r.Context(), pathValueInt)

		if err != nil {
			if err == sql.ErrNoRows {

				http.Error(w, "Product not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(product)
	}
}

func AddProductHandler(store ProductStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product *ProductInput

		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, "Malformed JSON", http.StatusBadRequest)
			return
		}

		if err := store.AddProduct(r.Context(), &Product{Name: product.Name, Price: product.Price, Quantity: product.Quantity, CategoryID: product.CategoryID}); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func RemoveProductHandler(store ProductStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathValue := r.PathValue("id")
		pathValueInt, err := strconv.Atoi(pathValue)

		if err != nil {
			http.Error(w, "Invalid path value", http.StatusBadRequest)
			return
		}

		if err := store.RemoveProduct(r.Context(), &Product{ID: pathValueInt}); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Product not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
func UpdateProductHandler(store ProductStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product *ProductInput

		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, "Malformed JSON", http.StatusBadRequest)
			return
		}

		pathValue := r.PathValue("id")
		pathValueInt, err := strconv.Atoi(pathValue)

		if err != nil {
			http.Error(w, "Invalid path value", http.StatusBadRequest)
			return
		}

		if err := store.UpdateProduct(r.Context(), &Product{ID: pathValueInt, Name: product.Name, Price: product.Price, Quantity: product.Quantity, CategoryID: product.CategoryID}); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Product not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
