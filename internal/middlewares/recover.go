package middlewares

import (
	"github.com/saurav11sarkar/go/internal/utils"
	"log"
	"net/http"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if err == http.ErrAbortHandler {
					panic(err)
				}
				log.Printf("request_id=%q panic recovered", r.Header.Get("X-Request-Id"))
				utils.HandlerError(w, utils.NewAppError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Internal server error"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
