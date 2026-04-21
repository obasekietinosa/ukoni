package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"ukoni/agent/internal/agent"
	"ukoni/agent/internal/config"
	"ukoni/agent/internal/database"
	"ukoni/agent/internal/handler"
)

func main() {
	cfg := config.Load()

	log.Printf("Starting agent server on port %d", cfg.Port)

	// Initialize Database
	dbService, err := database.New(cfg.DBURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbService.Close()
	log.Println("database connected")

	// Initialize Agent
	a := agent.New(cfg.APIBaseURL)

	h := &handler.AgentHandler{
		Agent: a,
	}

	mux := http.NewServeMux()
	// Using Go 1.22+ pattern matching for method
	mux.HandleFunc("POST /chat", h.HandleChat)

	if err := http.ListenAndServe(":"+strconv.Itoa(cfg.Port), enableCORS(mux)); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		allowed := false
		for _, o := range config.Load().CorsAllowedOrigins {
			o = strings.TrimSpace(o)
			o = strings.TrimRight(o, "/")
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")

			if reqHeaders := r.Header.Get("Access-Control-Request-Headers"); reqHeaders != "" {
				w.Header().Set("Access-Control-Allow-Headers", reqHeaders)
			} else {
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
		} else if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
