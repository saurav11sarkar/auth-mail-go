package routes

import (
	"net/http"

	"github.com/saurav11sarkar/go/internal/auth"
	"github.com/saurav11sarkar/go/internal/categories"
	"github.com/saurav11sarkar/go/internal/config"
	"github.com/saurav11sarkar/go/internal/middlewares"
	"github.com/saurav11sarkar/go/internal/user"
	"github.com/saurav11sarkar/go/internal/utils"
)

func NewRouter(deps *Deps, cfg config.Config) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		utils.JSON(w, http.StatusOK, "Server is running", nil)
	})

	v1Mux := http.NewServeMux()

	auth.RegisterRoutes(v1Mux, deps.Auth)
	user.RegisterRoutes(v1Mux, deps.User, cfg.Auth.JwtAccessSecret)
	categories.CategoryRouter(v1Mux, deps.Category, cfg.Auth.JwtAccessSecret)

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1Mux))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		utils.Error(w, http.StatusNotFound, "Route not found")
	})

	return middlewares.Chain(mux,
		middlewares.RequestID,
		middlewares.Logger,
		middlewares.Recover,
		middlewares.CORS(cfg.CorsOrigin),
	)
}
