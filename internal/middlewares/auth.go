package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/saurav11sarkar/go/internal/config"
	"github.com/saurav11sarkar/go/internal/utils"
)

type contextKey string

const UserContextKey contextKey = "user"

func Auth(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cfg := config.MustLoad()
			token := strings.TrimSpace(r.Header.Get("Authorization"))
			parts := strings.SplitN(token, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.HandlerError(w, utils.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized"))
				return
			}
			cliem, err := utils.ParseJwtToken(parts[1], cfg.Auth.JwtAccessSecret, "access")

			if err != nil {
				utils.HandlerError(w, utils.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized"))
				return
			}

			if len(roles) > 0 {
				hashRole := false
				for _, role := range roles {
					if role == cliem.Role {
						hashRole = true
						break
					}
				}
				if !hashRole {
					utils.HandlerError(w, utils.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized"))
					return
				}
			}

			ctx := context.WithValue(r.Context(), UserContextKey, cliem)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
