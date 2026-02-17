package handlers

import (
	"encoding/json"
	"net/http"
	"ukoni/internal/services"
)

type SellerHandler struct {
	Service *services.SellerService
}

type SellerRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// CreateSeller creates a new seller.
// @Summary Create seller
// @Description Create a new seller.
// @Tags Seller
// @Accept json
// @Produce json
// @Param request body SellerRequest true "Seller Request"
// @Success 201 {object} models.Seller
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /sellers [post]
// @Security BearerAuth
func (h *SellerHandler) CreateSeller(w http.ResponseWriter, r *http.Request) {
	var req SellerRequest
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

// GetSeller retrieves a seller by ID.
// @Summary Get seller
// @Description Get details of a specific seller.
// @Tags Seller
// @Produce json
// @Param id path string true "Seller ID"
// @Success 200 {object} models.Seller
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "seller not found"
// @Failure 500 {string} string "internal server error"
// @Router /sellers/{id} [get]
// @Security BearerAuth
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

// ListSellers lists all sellers.
// @Summary List sellers
// @Description Retrieve a list of all sellers.
// @Tags Seller
// @Produce json
// @Success 200 {array} models.Seller
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /sellers [get]
// @Security BearerAuth
func (h *SellerHandler) ListSellers(w http.ResponseWriter, r *http.Request) {
	sellers, err := h.Service.ListSellers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(sellers)
}

// UpdateSeller updates a seller.
// @Summary Update seller
// @Description Update details of an existing seller.
// @Tags Seller
// @Accept json
// @Produce json
// @Param id path string true "Seller ID"
// @Param request body SellerRequest true "Seller Request"
// @Success 200 {object} models.Seller
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "seller not found"
// @Failure 500 {string} string "internal server error"
// @Router /sellers/{id} [put]
// @Security BearerAuth
func (h *SellerHandler) UpdateSeller(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "seller id required", http.StatusBadRequest)
		return
	}

	var req SellerRequest
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

// DeleteSeller deletes a seller.
// @Summary Delete seller
// @Description Delete a seller.
// @Tags Seller
// @Param id path string true "Seller ID"
// @Success 204
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "seller not found"
// @Failure 500 {string} string "internal server error"
// @Router /sellers/{id} [delete]
// @Security BearerAuth
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
