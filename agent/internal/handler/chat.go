package handler

import (
	"encoding/json"
	"net/http"
	"ukoni/agent/internal/agent"
)

type AgentHandler struct {
	Agent *agent.Agent
}

type ChatRequest struct {
	Prompt      string `json:"prompt"`
	InventoryID string `json:"inventory_id"`
}

type ChatResponse struct {
	Response string `json:"response"`
	Error    string `json:"error,omitempty"`
}

func (h *AgentHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.InventoryID == "" {
		http.Error(w, "inventory_id required", http.StatusBadRequest)
		return
	}

	// Extract Auth Token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "missing authorization header", http.StatusUnauthorized)
		return
	}

	respStr, err := h.Agent.Run(r.Context(), req.Prompt, req.InventoryID, authHeader)

	resp := ChatResponse{
		Response: respStr,
	}
	if err != nil {
		resp.Error = err.Error()
		// If error is related to settings missing, maybe 400?
		// For now 500 is safe default for server-side processing errors.
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
