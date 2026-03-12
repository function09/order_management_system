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

type LoginInput struct {
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

		var userRegister User
		userRegister.PasswordHash = hashedPassword
		userRegister.Username = input.Username
		userRegister.CreatedAt = time.Now()

		if err := store.RegisterUser(&userRegister); err != nil {
			http.Error(w, "Error creating new user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}

}

func LoginUserHandler(store AuthStore, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var input LoginInput

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		user, err := store.GetUserByUserName(input.Username)

		if err != nil {
			http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
			return
		}

		err = VerifyPassword(user.PasswordHash, input.Password)

		if err != nil {
			http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
			return
		}

		token, err := GenerateToken(user.Username, secret, time.Hour)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			HttpOnly: true,
			Path:     "/",
		})

		w.WriteHeader(http.StatusOK)
	}
}

func LogOutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "token", Expires: time.Unix(0, 0)})
	w.WriteHeader(http.StatusOK)

}
