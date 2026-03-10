package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type FakeStore struct{}

func (f *FakeStore) RegisterUser(user *User) error {
	return nil
}

func (f *FakeStore) GetUserByUserName(username string) (*User, error) {
	return &User{}, nil
}

func TestRegisterUser(t *testing.T) {

	t.Run("Generates a new user", func(t *testing.T) {

		body := strings.NewReader(`{"username":"testuser", "password":"testpass"}`)
		req := httptest.NewRequest("POST", "/auth/register", body)
		w := httptest.NewRecorder()

		handler := RegisterUserHandler(&FakeStore{})

		handler(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("got %d want %d", w.Code, http.StatusCreated)
		}
	})

}
