package middlewares

import (
	"github.com/saurav11sarkar/go/internal/utils"
	"net/http"
)

// Method guards an exact route while keeping API errors in JSON.
func Method(method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method && !(method == http.MethodGet && r.Method == http.MethodHead) {
			w.Header().Set("Allow", method)
			if method == http.MethodGet {
				w.Header().Set("Allow", "GET, HEAD")
			}
			utils.HandlerError(w, utils.NewAppError(http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
