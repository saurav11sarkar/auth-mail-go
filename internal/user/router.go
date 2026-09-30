package user

import (
	"net/http"

	"github.com/saurav11sarkar/go/internal/middlewares"
)

func UserRouter(mux *http.ServeMux, h *Handler) {
	mux.Handle("GET /user/profile", middlewares.Auth("admin", "user")(http.HandlerFunc(h.Profile)))
}
