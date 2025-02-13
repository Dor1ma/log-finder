package middleware

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

func RateLimit(next http.HandlerFunc, requestsPerSecond int) http.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond)

	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		duration := time.Since(start)
		log.Printf("%s %s %s %v", r.Method, r.URL.Path, r.RemoteAddr, duration)
	}
}
