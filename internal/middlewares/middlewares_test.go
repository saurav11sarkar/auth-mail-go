package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/saurav11sarkar/go/internal/utils"
)

func TestAuthInjectedSecretAndRoles(t *testing.T) {
	secret := "test-secret"
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(UserContextKey).(*utils.JwtClaims)
		if !ok || claims.UserID != "user-id" {
			t.Fatal("missing claims")
		}
		w.WriteHeader(204)
	})
	for _, tc := range []struct {
		role   string
		status int
	}{{"user", 204}, {"admin", 403}} {
		token, err := utils.CreateJwtToken("user-id", tc.role, "access", secret, time.Minute)
		if err != nil {
			t.Fatal(err)
		}
		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set("Authorization", "bearer "+token)
		w := httptest.NewRecorder()
		Auth(secret, "user")(next).ServeHTTP(w, r)
		if w.Code != tc.status {
			t.Fatalf("got %d", w.Code)
		}
	}
}

func TestChainOrderAndRecovery(t *testing.T) {
	var order []string
	wrap := func(name string) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, name+" before")
				next.ServeHTTP(w, r)
				order = append(order, name+" after")
			})
		}
	}
	w := httptest.NewRecorder()
	Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { panic("secret") }), wrap("outer"), wrap("inner"), Recover).ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if strings.Join(order, ",") != "outer before,inner before,inner after,outer after" {
		t.Fatal(order)
	}
	if w.Code != 500 || strings.Contains(w.Body.String(), "secret") || w.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("bad recovery: %d %s", w.Code, w.Body.String())
	}
}
