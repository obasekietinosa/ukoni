package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"ukoni/internal/config"
	"ukoni/internal/handlers"
)

func TestAgentHandler_HandleProxy(t *testing.T) {
	// Start a mock backend server (Agent Service)
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify path
		if r.URL.Path != "/chat" {
			t.Errorf("expected path /chat, got %s", r.URL.Path)
		}
		// Verify method
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}
		// Verify header
		if r.Header.Get("Authorization") != "Bearer token123" {
			t.Errorf("expected Authorization header, got %s", r.Header.Get("Authorization"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("agent response"))
	}))
	defer backendServer.Close()

	// Configure handler
	cfg := &config.Config{
		AgentServiceURL: backendServer.URL,
	}
	agentHandler, err := handlers.NewAgentHandler(cfg)
	if err != nil {
		t.Fatalf("failed to create agent handler: %v", err)
	}

	// Create a request
	req, err := http.NewRequest("POST", "/agent/chat", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer token123")

	// Create a recorder
	rr := httptest.NewRecorder()

	// Handler with StripPrefix, mimicking the router setup
	handler := http.StripPrefix("/agent", http.HandlerFunc(agentHandler.HandleProxy))

	// Serve
	handler.ServeHTTP(rr, req)

	// Check response
	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", rr.Code)
	}
	if rr.Body.String() != "agent response" {
		t.Errorf("expected body 'agent response', got %s", rr.Body.String())
	}
}
