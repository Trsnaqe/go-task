package middlewares

import (
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

func RateLimitMiddleware(limiter *rate.Limiter) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
