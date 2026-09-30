package routes

import (
	"github.com/saurav11sarkar/go/internal/auth"
	"github.com/saurav11sarkar/go/internal/user"
)

type Deps struct {
	Auth *auth.Handler
	User *user.Handler
}
