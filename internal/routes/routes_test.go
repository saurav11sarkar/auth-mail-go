package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/saurav11sarkar/go/internal/auth"
	"github.com/saurav11sarkar/go/internal/config"
	"github.com/saurav11sarkar/go/internal/user"
)

// Exercise actual API mounting and middleware without a live database.
func TestAPIRouting(t *testing.T) {
	cfg := config.Config{CorsOrigin: "http://localhost:3000"}
	cfg.Auth.JwtAccessSecret = "test-secret"
	router := NewRouter(&Deps{Auth: auth.NewHandler(nil, cfg), User: user.NewHandler(nil)}, cfg)
	cases := []struct {
		method, path, body string
		status             int
	}{
		{"GET", "/", "", 200},
		{"GET", "/api/v1/auth/login", "", 405},
		{"POST", "/api/v1/user/profile", "", 405},
		{"GET", "/api/v1/user/profile", "", 401},
		{"POST", "/api/v1/auth/register", "{}", 400},
		{"POST", "/api/v1/auth/login", "{", 400},
		{"POST", "/api/v1/auth/refresh", "{}", 400},
		{"POST", "/api/v1/auth/forget-password", "{}", 400},
		{"POST", "/api/v1/auth/reset-password", "{}", 400},
		{"GET", "/api/v1/missing", "", 404},
		{"OPTIONS", "/api/v1/auth/login", "", 204},
	}
	for _, tc := range cases {
		t.Run(tc.method+tc.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body)))
			if w.Code != tc.status {
				t.Fatalf("got %d: %s", w.Code, w.Body.String())
			}
			if w.Header().Get("X-Request-Id") == "" {
				t.Fatal("missing request ID")
			}
			if tc.status != http.StatusNoContent && w.Header().Get("Content-Type") != "application/json" {
				t.Fatal("expected JSON")
			}
			if tc.status == 405 && w.Header().Get("Allow") == "" {
				t.Fatal("missing Allow header")
			}
		})
	}
}
