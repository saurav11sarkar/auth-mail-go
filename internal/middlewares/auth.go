package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/saurav11sarkar/go/internal/utils"
)

type contextKey string

const UserContextKey contextKey = "user"

func Auth(accessSecret string, roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := strings.TrimSpace(r.Header.Get("Authorization"))
			parts := strings.Fields(token)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				utils.HandlerError(w, utils.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized"))
				return
			}
			claims, err := utils.ParseJwtToken(parts[1], accessSecret, "access")

			if err != nil {
				utils.HandlerError(w, utils.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized"))
				return
			}

			if len(roles) > 0 {
				hasRole := false
				for _, role := range roles {
					if role == claims.Role {
						hasRole = true
						break
					}
				}
				if !hasRole {
					utils.HandlerError(w, utils.NewAppError(http.StatusForbidden, "FORBIDDEN", "Forbidden"))
					return
				}
			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
