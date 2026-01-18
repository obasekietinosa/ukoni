package handlers

import (
	"encoding/json"
	"net/http"
	"ukoni/internal/services"
)

type InventoryHandler struct {
	Service *services.InventoryService
}

func (h *InventoryHandler) CreateInventory(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	inventory, err := h.Service.CreateInventory(userID, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(inventory)
}

func (h *InventoryHandler) GetInventory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "inventory id required", http.StatusBadRequest)
		return
	}

	inventory, err := h.Service.GetInventory(id)
	if err != nil {
		// Differentiate between 404 and 500 if possible, but basic 500 for now or checks in service.
		// For simplicity/robustness assuming if error is not nil it might be not found or db error.
		// Ideally service returns named errors.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(inventory)
}

func (h *InventoryHandler) ListInventories(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	inventories, err := h.Service.ListInventories(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(inventories)
}
