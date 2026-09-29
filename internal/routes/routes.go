package routes

import (
	"net/http"

	"github.com/saurav11sarkar/go/internal/auth"
	"github.com/saurav11sarkar/go/internal/config"
	"github.com/saurav11sarkar/go/internal/middlewares"
	"github.com/saurav11sarkar/go/internal/utils"
)

func NewRouter(deps *Deps, cfg config.Config) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		utils.JSON(w, http.StatusOK, "Server is running", nil)
	})

	v1Mux := http.NewServeMux()


	auth.AuthRouter(v1Mux, deps.Auth)


	v1Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		utils.Error(w, http.StatusNotFound, "API v1 endpoint not found")
	})
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1Mux))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		utils.Error(w, http.StatusNotFound, "Route not found")
	})

	return middlewares.Chain(mux,
		middlewares.CORS(cfg.CorsOrigin),
		middlewares.RequestID,
		middlewares.Logger,
		middlewares.Recover,
	)
}
