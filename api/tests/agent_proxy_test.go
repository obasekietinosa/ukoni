package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgentProxy(t *testing.T) {
	clearDB()

	// Start a mock agent server
	mockAgent := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify path
		if r.URL.Path != "/chat" {
			t.Errorf("Expected path /chat, got %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Verify method
		if r.Method != "POST" {
			t.Errorf("Expected method POST, got %s", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Verify Auth header is passed through
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			t.Errorf("Expected Authorization header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Respond with success
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response": "Hello from mock agent"}`))
	}))
	defer mockAgent.Close()

	// Update config to point to mock agent
	// cfg is global in setup_test.go
	originalAgentURL := cfg.AgentServiceURL
	cfg.AgentServiceURL = mockAgent.URL
	defer func() { cfg.AgentServiceURL = originalAgentURL }()

	router := setupRouter()

	t.Run("Successful Proxy", func(t *testing.T) {
		// Create a user and get token
		token := createTestUser(router)

		payload := map[string]string{
			"prompt":       "Hello",
			"inventory_id": "test-inv-id",
		}
		body, _ := json.Marshal(payload)

		// Request to API server at /agent/chat
		// The test router (mux) matches /agent/ prefix
		req, _ := http.NewRequest("POST", "/agent/chat", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]string
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Hello from mock agent", response["response"])
	})

	t.Run("Proxy Unauthorized", func(t *testing.T) {
		payload := map[string]string{
			"prompt":       "Hello",
			"inventory_id": "test-inv-id",
		}
		body, _ := json.Marshal(payload)

		// Request without token
		req, _ := http.NewRequest("POST", "/agent/chat", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Should be blocked by Auth middleware on API server
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}
