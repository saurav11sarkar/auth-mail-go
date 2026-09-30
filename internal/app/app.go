package app

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/saurav11sarkar/go/internal/auth"
	"github.com/saurav11sarkar/go/internal/config"
	"github.com/saurav11sarkar/go/internal/routes"
	"github.com/saurav11sarkar/go/internal/user"
)

func NewHandler(db *pgxpool.Pool, cfg config.Config) (http.Handler, error) {

	authrepo := auth.NewResposistory(db)
	authservice := auth.NewService(authrepo, cfg)
	authhandler := auth.NewHandler(authservice)

	userrepo := user.NewRepository(db)
	userservice := user.NewService(userrepo, cfg)
	userhandler := user.NewHandler(userservice)

	deps := routes.Deps{
		Auth: authhandler,
		User: userhandler,
	}

	router := routes.NewRouter(&deps, cfg)
	return router, nil
}
