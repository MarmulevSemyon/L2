package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware логирует метод, URL и время обработки каждого HTTP-запроса.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Printf("method: %s\turl: %s\ttime: %s\n", r.Method, r.URL.String(), duration)
	})
}
