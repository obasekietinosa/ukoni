package middleware

import (
	"net/http"
	"strings"
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
			o = strings.TrimSpace(o)
			o = strings.TrimRight(o, "/")
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			// Echo origin instead of wildcard for broader compatibility, especially with credentials
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")

			// Dynamically echo requested headers if present, else default to common ones
			if reqHeaders := r.Header.Get("Access-Control-Request-Headers"); reqHeaders != "" {
				w.Header().Set("Access-Control-Allow-Headers", reqHeaders)
			} else {
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
			}

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
		} else if r.Method == http.MethodOptions {
			// Preflight requests should not fall through to the mux if not explicitly allowed,
			// otherwise the mux may return a 405 Method Not Allowed.
			// Return a 204 No Content to cleanly trigger the browser's CORS failure.
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
