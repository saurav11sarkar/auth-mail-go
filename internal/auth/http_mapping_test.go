package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/saurav11sarkar/go/internal/config"
	"github.com/saurav11sarkar/go/internal/utils"
)

func TestRefreshCookiePolicy(t *testing.T) {
	for _, env := range []string{"production", "development"} {
		cfg := config.Config{Env: env}
		cfg.Auth.RefreshTokenDays = 7
		h := NewHandler(nil, cfg)
		w := httptest.NewRecorder()
		h.setRefreshCookie(w, "token")
		cookies := w.Result().Cookies()
		if len(cookies) != 1 {
			t.Fatalf("expected one cookie, got %d", len(cookies))
		}
		c := cookies[0]
		if c.Name != "refresh_token" || c.Value != "token" || c.Path != "/api/v1/auth" || c.MaxAge != 604800 || !c.HttpOnly || c.Secure != (env == "production") {
			t.Fatalf("incorrect cookie for %s: %+v", env, c)
		}
	}
}

func TestPublicResponsesPreserveContractAndHideCredentials(t *testing.T) {
	now := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	account := Account{ID: "id", Email: "a@example.com", Name: "Alice", Role: "user", Status: "active", Password: "secret-hash", OTP: "123456", OTPExpiresAt: now, CreatedAt: now, UpdatedAt: now}
	cases := []struct {
		name  string
		value any
		want  string
	}{
		{"register", registerResponse(account), `{"id":"id","email":"a@example.com","name":"Alice","role":"user","status":"active","otp_expires_at":"0001-01-01T00:00:00Z","created_at":"2026-10-01T00:00:00Z","updated_at":"2026-10-01T00:00:00Z"}`},
		{"login", loginResponse(LoginResult{Account: account, AccessToken: "access", RefreshToken: "refresh"}), `{"id":"id","email":"a@example.com","name":"Alice","role":"user","status":"active","access_token":"access","refresh_token":"refresh","created_at":"2026-10-01T00:00:00Z","updated_at":"2026-10-01T00:00:00Z"}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.value)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tc.want {
				t.Fatalf("JSON contract changed:\ngot  %s\nwant %s", got, tc.want)
			}
		})
	}
}

func TestHTTPErrorMapping(t *testing.T) {
	cases := []struct {
		err    error
		status int
		code   string
	}{
		{fmt.Errorf("wrapped: %w", ErrEmailExists), 409, "EMAIL_ALREADY_EXISTS"},
		{ErrInvalidCredentials, 401, "UNAUTHORIZED"},
		{ErrInvalidRefreshToken, 401, "UNAUTHORIZED"},
		{ErrInvalidOTP, 400, "INVALID_OR_EXPIRED_OTP"},
		{utils.NewAppError(400, "VALIDATION_ERROR", "Invalid input"), 400, "VALIDATION_ERROR"},
		{errors.New("private database detail"), 500, "INTERNAL_SERVER_ERROR"},
	}
	for _, tc := range cases {
		w := httptest.NewRecorder()
		writeError(w, tc.err)
		if w.Code != tc.status || !strings.Contains(w.Body.String(), tc.code) {
			t.Fatalf("got %d %s", w.Code, w.Body.String())
		}
		if strings.Contains(w.Body.String(), "private database detail") {
			t.Fatal("internal error leaked")
		}
	}
}

func TestRepositoryErrorTranslation(t *testing.T) {
	duplicate := &pgconn.PgError{Code: "23505", ConstraintName: "users_email_key"}
	if !errors.Is(accountRepositoryError(duplicate), ErrEmailExists) {
		t.Fatal("duplicate email not translated")
	}
	other := &pgconn.PgError{Code: "23505", ConstraintName: "users_pkey"}
	if errors.Is(accountRepositoryError(other), ErrEmailExists) {
		t.Fatal("unrelated constraint treated as email conflict")
	}
	if !errors.Is(accountRepositoryError(pgx.ErrNoRows), ErrAccountNotFound) {
		t.Fatal("missing account not translated")
	}
	cause := errors.New("connection failed")
	if !errors.Is(accountRepositoryError(cause), cause) {
		t.Fatal("underlying error lost")
	}
}
