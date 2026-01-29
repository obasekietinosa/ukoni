package middleware

import (
	"net/http"
	"ukoni/internal/config"
)

type CorsMiddleware struct {
	Config *config.Config
}

func NewCorsMiddleware(cfg *config.Config) *CorsMiddleware {
	return &CorsMiddleware{Config: cfg}
}

func (m *CorsMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// If no Origin header is present, just pass the request to the next handler.
		// It's likely a same-origin request or a non-browser client.
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		allowed := false
		for _, o := range m.Config.CorsAllowedOrigins {
			if o == "*" {
				allowed = true
				w.Header().Set("Access-Control-Allow-Origin", "*")
				break
			}
			if o == origin {
				allowed = true
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		if allowed {
			w.Header().Add("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
