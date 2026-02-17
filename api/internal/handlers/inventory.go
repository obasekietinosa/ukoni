package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"ukoni/internal/services"
)

type InventoryHandler struct {
	Service                 *services.InventoryService
	InventoryProductService *services.InventoryProductService
}

type CreateInventoryRequest struct {
	Name string `json:"name"`
}

// CreateInventory creates a new inventory.
// @Summary Create inventory
// @Description Create a new inventory for the user.
// @Tags Inventory
// @Accept json
// @Produce json
// @Param request body CreateInventoryRequest true "Create Inventory Request"
// @Success 201 {object} models.Inventory
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /inventories [post]
// @Security BearerAuth
func (h *InventoryHandler) CreateInventory(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateInventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	inventory, err := h.Service.CreateInventory(r.Context(), userID, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(inventory)
}

// GetInventory retrieves an inventory by ID.
// @Summary Get inventory
// @Description Get details of a specific inventory.
// @Tags Inventory
// @Produce json
// @Param id path string true "Inventory ID"
// @Success 200 {object} models.Inventory
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /inventories/{id} [get]
// @Security BearerAuth
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

// ListInventories lists all inventories for the user.
// @Summary List inventories
// @Description Retrieve a list of inventories the user belongs to.
// @Tags Inventory
// @Produce json
// @Success 200 {array} models.Inventory
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /inventories [get]
// @Security BearerAuth
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventories)
}

// ListInventoryProducts lists all products in an inventory.
// @Summary List inventory products
// @Description Retrieve a list of products in the specified inventory.
// @Tags Inventory
// @Produce json
// @Param id path string true "Inventory ID"
// @Success 200 {array} models.InventoryProduct
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 500 {string} string "internal server error"
// @Router /inventories/{id}/inventory-products [get]
// @Security BearerAuth
func (h *InventoryHandler) ListInventoryProducts(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	inventoryID := r.PathValue("id")
	if inventoryID == "" {
		http.Error(w, "inventory id required", http.StatusBadRequest)
		return
	}

	// Verify membership
	_, err := h.Service.MembershipModel.GetMembership(inventoryID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "forbidden", http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	products, err := h.InventoryProductService.ListInventoryProducts(r.Context(), inventoryID, 100, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
