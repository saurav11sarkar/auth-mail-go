package user

import (
	"net/http"

	"github.com/saurav11sarkar/go/internal/middlewares"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler, accessSecret string) {
	profile := middlewares.Auth(accessSecret, "admin", "user")(http.HandlerFunc(h.Profile))
	mux.Handle("GET /user/profile", middlewares.Method(http.MethodGet, profile))
	updateProfile := middlewares.Auth(accessSecret, "admin", "user")(http.HandlerFunc(h.UpdateProfile))
	mux.Handle("PUT /user/profile", middlewares.Method(http.MethodPut, updateProfile))
}
