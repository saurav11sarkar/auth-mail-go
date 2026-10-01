package auth

import (
	"github.com/saurav11sarkar/go/internal/middlewares"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.Handle("/auth/register", middlewares.Method(http.MethodPost, http.HandlerFunc(h.Create)))
	mux.Handle("/auth/login", middlewares.Method(http.MethodPost, http.HandlerFunc(h.Login)))
	mux.Handle("/auth/refresh", middlewares.Method(http.MethodPost, http.HandlerFunc(h.Refresh)))
	mux.Handle("/auth/forget-password", middlewares.Method(http.MethodPost, http.HandlerFunc(h.ForgetPassword)))
	mux.Handle("/auth/reset-password", middlewares.Method(http.MethodPost, http.HandlerFunc(h.ResetPassword)))
}
