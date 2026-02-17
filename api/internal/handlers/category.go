package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"ukoni/internal/services"
)

type CategoryHandler struct {
	Service           *services.CategoryService
	MembershipService *services.MembershipService
}

type CategoryRequest struct {
	Name             string  `json:"name"`
	ParentCategoryID *string `json:"parent_category_id"`
}

// CreateCategory creates a new category.
// @Summary Create category
// @Description Create a new category in the inventory.
// @Tags Category
// @Accept json
// @Produce json
// @Param id path string true "Inventory ID"
// @Param request body CategoryRequest true "Category Request"
// @Success 201 {object} models.Category
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 500 {string} string "internal server error"
// @Router /inventories/{id}/categories [post]
// @Security BearerAuth
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	inventoryID := r.PathValue("id")
	if inventoryID == "" {
		http.Error(w, "inventory id required", http.StatusBadRequest)
		return
	}

	// Check membership
	if _, err := h.MembershipService.MembershipModel.GetMembership(inventoryID, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "forbidden", http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var req CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	category, err := h.Service.CreateCategory(r.Context(), inventoryID, req.Name, req.ParentCategoryID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

// ListCategories lists all categories in an inventory.
// @Summary List categories
// @Description Retrieve a list of categories in the specified inventory.
// @Tags Category
// @Produce json
// @Param id path string true "Inventory ID"
// @Success 200 {array} models.Category
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 500 {string} string "internal server error"
// @Router /inventories/{id}/categories [get]
// @Security BearerAuth
func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	inventoryID := r.PathValue("id")
	if inventoryID == "" {
		http.Error(w, "inventory id required", http.StatusBadRequest)
		return
	}

	// Check membership
	if _, err := h.MembershipService.MembershipModel.GetMembership(inventoryID, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "forbidden", http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	categories, err := h.Service.ListCategories(r.Context(), inventoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
