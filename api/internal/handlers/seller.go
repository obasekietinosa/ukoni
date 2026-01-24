package handlers

import (
	"encoding/json"
	"net/http"
	"ukoni/internal/services"
)

type SellerHandler struct {
	Service *services.SellerService
}

func (h *SellerHandler) CreateSeller(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	seller, err := h.Service.CreateSeller(req.Name, req.Type)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(seller)
}

func (h *SellerHandler) GetSeller(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "seller id required", http.StatusBadRequest)
		return
	}

	seller, err := h.Service.GetSeller(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(seller)
}

func (h *SellerHandler) ListSellers(w http.ResponseWriter, r *http.Request) {
	sellers, err := h.Service.ListSellers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(sellers)
}

func (h *SellerHandler) UpdateSeller(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "seller id required", http.StatusBadRequest)
		return
	}

	var req struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	seller, err := h.Service.UpdateSeller(id, req.Name, req.Type)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(seller)
}

func (h *SellerHandler) DeleteSeller(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "seller id required", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteSeller(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
