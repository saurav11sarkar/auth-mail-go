package auth

import "net/http"

func AuthRouter(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/auth/register", h.Create)
	mux.HandleFunc("/auth/login", h.Login)
	mux.HandleFunc("/auth/refresh", h.Refresh)
	mux.HandleFunc("/auth/forget-password", h.ForgetPassword)
	mux.HandleFunc("/auth/reset-password", h.ResetPassword)
}
