package categories

import (
	"net/http"

	"github.com/saurav11sarkar/go/internal/middlewares"
)

func CategoryRouter(mux *http.ServeMux, h *Handler, accessSecret string) {
	mux.Handle("PUT /categories/{id}", middlewares.Auth(accessSecret, "admin")(http.HandlerFunc(h.Update)))
	mux.Handle("DELETE /categories/{id}", middlewares.Auth(accessSecret, "admin")(http.HandlerFunc(h.Delete)))
	mux.Handle("POST /categories", middlewares.Auth(accessSecret, "admin")(http.HandlerFunc(h.Create)))
	mux.Handle("GET /categories", middlewares.Auth(accessSecret, "admin")(http.HandlerFunc(h.GetAllCategories)))
	mux.HandleFunc("GET /categories/{id}", h.GetCategoryByID)
}
