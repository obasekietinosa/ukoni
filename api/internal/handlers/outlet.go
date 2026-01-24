package handlers

import (
	"encoding/json"
	"net/http"
	"ukoni/internal/services"
)

type OutletHandler struct {
	Service *services.OutletService
}

func (h *OutletHandler) CreateOutlet(w http.ResponseWriter, r *http.Request) {
	sellerID := r.PathValue("id")
	if sellerID == "" {
		http.Error(w, "seller id required", http.StatusBadRequest)
		return
	}

	var req struct {
		Name       string `json:"name"`
		Channel    string `json:"channel"`
		Address    string `json:"address"`
		WebsiteURL string `json:"website_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	outlet, err := h.Service.CreateOutlet(sellerID, req.Name, req.Channel, req.Address, req.WebsiteURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(outlet)
}

func (h *OutletHandler) GetOutlet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "outlet id required", http.StatusBadRequest)
		return
	}

	outlet, err := h.Service.GetOutlet(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(outlet)
}

func (h *OutletHandler) ListOutlets(w http.ResponseWriter, r *http.Request) {
	sellerID := r.PathValue("id")
	if sellerID == "" {
		http.Error(w, "seller id required", http.StatusBadRequest)
		return
	}

	outlets, err := h.Service.ListOutlets(sellerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(outlets)
}

func (h *OutletHandler) UpdateOutlet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "outlet id required", http.StatusBadRequest)
		return
	}

	var req struct {
		Name       string `json:"name"`
		Channel    string `json:"channel"`
		Address    string `json:"address"`
		WebsiteURL string `json:"website_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	outlet, err := h.Service.UpdateOutlet(id, req.Name, req.Channel, req.Address, req.WebsiteURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(outlet)
}

func (h *OutletHandler) DeleteOutlet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "outlet id required", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteOutlet(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
