package handlers

import (
	"encoding/json"
	"net/http"
	"ukoni/internal/services"
)

type OutletHandler struct {
	Service *services.OutletService
}

type OutletRequest struct {
	Name       string `json:"name"`
	Channel    string `json:"channel"`
	Address    string `json:"address"`
	WebsiteURL string `json:"website_url"`
}

// CreateOutlet creates a new outlet.
// @Summary Create outlet
// @Description Create a new outlet for a seller.
// @Tags Outlet
// @Accept json
// @Produce json
// @Param id path string true "Seller ID"
// @Param request body OutletRequest true "Outlet Request"
// @Success 201 {object} models.Outlet
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "seller not found"
// @Failure 500 {string} string "internal server error"
// @Router /sellers/{id}/outlets [post]
// @Security BearerAuth
func (h *OutletHandler) CreateOutlet(w http.ResponseWriter, r *http.Request) {
	sellerID := r.PathValue("id")
	if sellerID == "" {
		http.Error(w, "seller id required", http.StatusBadRequest)
		return
	}

	var req OutletRequest
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

// GetOutlet retrieves an outlet by ID.
// @Summary Get outlet
// @Description Get details of a specific outlet.
// @Tags Outlet
// @Produce json
// @Param id path string true "Outlet ID"
// @Success 200 {object} models.Outlet
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "outlet not found"
// @Failure 500 {string} string "internal server error"
// @Router /outlets/{id} [get]
// @Security BearerAuth
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

// ListOutlets lists outlets for a seller.
// @Summary List outlets
// @Description Retrieve a list of outlets for a seller.
// @Tags Outlet
// @Produce json
// @Param id path string true "Seller ID"
// @Success 200 {array} models.Outlet
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "seller not found"
// @Failure 500 {string} string "internal server error"
// @Router /sellers/{id}/outlets [get]
// @Security BearerAuth
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

// UpdateOutlet updates an outlet.
// @Summary Update outlet
// @Description Update details of an existing outlet.
// @Tags Outlet
// @Accept json
// @Produce json
// @Param id path string true "Outlet ID"
// @Param request body OutletRequest true "Outlet Request"
// @Success 200 {object} models.Outlet
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "outlet not found"
// @Failure 500 {string} string "internal server error"
// @Router /outlets/{id} [put]
// @Security BearerAuth
func (h *OutletHandler) UpdateOutlet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "outlet id required", http.StatusBadRequest)
		return
	}

	var req OutletRequest
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

// DeleteOutlet deletes an outlet.
// @Summary Delete outlet
// @Description Delete an outlet.
// @Tags Outlet
// @Param id path string true "Outlet ID"
// @Success 204
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "outlet not found"
// @Failure 500 {string} string "internal server error"
// @Router /outlets/{id} [delete]
// @Security BearerAuth
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
