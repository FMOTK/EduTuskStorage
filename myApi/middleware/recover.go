package middleware

import (
	"log"
	"net/http"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if rec := recover(); rec != nil {
				http.Error(w, "Internal server Error", http.StatusInternalServerError)
				log.Printf("PANIC RECOVERED: %v", rec)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
