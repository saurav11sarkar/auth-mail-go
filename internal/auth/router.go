package auth

import (
	"net/http"

	"github.com/saurav11sarkar/go/internal/middlewares"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.Handle("POST /auth/register", middlewares.Method(http.MethodPost, http.HandlerFunc(h.Create)))
	mux.Handle("POST /auth/login", middlewares.Method(http.MethodPost, http.HandlerFunc(h.Login)))
	mux.Handle("POST /auth/refresh", middlewares.Method(http.MethodPost, http.HandlerFunc(h.Refresh)))
	mux.Handle("POST /auth/forget-password", middlewares.Method(http.MethodPost, http.HandlerFunc(h.ForgetPassword)))
	mux.Handle("POST /auth/reset-password", middlewares.Method(http.MethodPost, http.HandlerFunc(h.ResetPassword)))
}
