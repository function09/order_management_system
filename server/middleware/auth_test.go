package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/function09/order_management_system/server/internal/auth"
)

func TestAuthMiddleware(t *testing.T) {
	token, _ := auth.GenerateToken("username", "secret", 1*time.Hour)
	var tests = []struct {
		name    string
		cookies bool
		want    int
	}{
		{"Unauthorized if no session", false, 401},
		{"Authorized if session", true, 200},
	}
	for _, e := range tests {
		t.Run(e.name, func(t *testing.T) {

			req := httptest.NewRequest("GET", "/products", nil)
			w := httptest.NewRecorder()

			dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			if e.cookies {
				req.AddCookie(&http.Cookie{
					Name:  "token",
					Value: token,
				})
			}

			protectedHandler := AuthMiddleware("secret", dummyHandler)
			protectedHandler(w, req)

			if w.Code != e.want {
				t.Errorf("Got %d, want %d", w.Code, e.want)
			}
		})
	}
}
