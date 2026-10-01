package middlewares

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("request_id=%q method=%s path=%q duration=%s", r.Header.Get("X-Request-Id"), r.Method, r.URL.Path, time.Since(start))
	})
}
