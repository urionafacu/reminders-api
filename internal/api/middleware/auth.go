package middleware

import (
	"net/http"
)

func APIKeyAuth(validAPIKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				http.Error(w, "API key is missing", http.StatusUnauthorized)
				return
			}

			if apiKey != validAPIKey {
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				return
			}

			// API key is valid, call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
