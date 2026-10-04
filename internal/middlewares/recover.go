package middlewares

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/saurav11sarkar/go/internal/utils"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if err == http.ErrAbortHandler {
					panic(err)
				}
				log.Printf("[PANIC RECOVERED] request_id=%q error=%v\nstack:\n%s",
					r.Header.Get("X-Request-Id"),
					err,
					debug.Stack(),
				)
				utils.HandlerError(w, utils.NewAppError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Internal server error"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
