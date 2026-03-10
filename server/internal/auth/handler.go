package auth

import (
	"encoding/json"
	"net/http"
	"time"
)

type RegisterInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterUserHandler(store AuthStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var input RegisterInput

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		hashedPassword, err := HashPassword(input.Password)

		if err != nil {
			http.Error(w, "Error creating new user", http.StatusInternalServerError)
			return
		}

		var user User
		user.PasswordHash = hashedPassword
		user.Username = input.Username
		user.CreatedAt = time.Now()
		if err := store.RegisterUser(&user); err != nil {
			http.Error(w, "Error creating new user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}

}
