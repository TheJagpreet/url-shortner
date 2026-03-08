package middleware

import (
	"net/http"
	"os"
)

// APIKeyAuth is a middleware that checks for a valid API key in the request header.
func APIKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		validKey := os.Getenv("URL_SHORTENER_API_KEY")
		if validKey == "" {
			// Fallback: allow if no key is set (for dev)
			http.Error(w, "API key not configured", http.StatusInternalServerError)
			return
		}
		if apiKey != validKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Proceed to next handler
		next.ServeHTTP(w, r)
	})
}
