package main

import (
	"log"
	"net/http"
	"strconv"
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
	mux.HandleFunc("POST /chat", enableCORS(h.HandleChat))
	// Handle OPTIONS explicitly or via middleware catch-all
	mux.HandleFunc("OPTIONS /chat", enableCORS(h.HandleChat)) // Just to handle OPTIONS

	if err := http.ListenAndServe(":"+strconv.Itoa(cfg.Port), mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Simple CORS handling
		// Ideally verify origin against config, but for now just reflect or allow specific
		// config.Load().CORSAllowed is used here
		w.Header().Set("Access-Control-Allow-Origin", config.Load().CORSAllowed)
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
