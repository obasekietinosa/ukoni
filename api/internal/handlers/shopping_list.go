package handlers

import (
	"encoding/json"
	"net/http"
	"ukoni/internal/models"
	"ukoni/internal/services"
)

type ShoppingListHandler struct {
	Service *services.ShoppingListService
}

type CreateListRequest struct {
	Name  string `json:"name"`
	Notes string `json:"notes"`
}

type AddItemRequest struct {
	TargetType  string   `json:"target_type"`
	TargetID    string   `json:"target_id"`
	Quantity    *float64 `json:"quantity"`
	Unit        string   `json:"unit"`
	Notes       string   `json:"notes"`
	ManualOrder *int     `json:"manual_order"`
}

type UpdateItemRequest struct {
	Quantity    *float64 `json:"quantity"`
	Unit        string   `json:"unit"`
	Notes       string   `json:"notes"`
	ManualOrder *int     `json:"manual_order"`
}

// CreateList creates a new shopping list.
// @Summary Create shopping list
// @Description Create a new shopping list for an inventory.
// @Tags Shopping List
// @Accept json
// @Produce json
// @Param id path string true "Inventory ID"
// @Param request body CreateListRequest true "Create List Request"
// @Success 201 {object} models.ShoppingList
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 500 {string} string "internal server error"
// @Router /inventories/{id}/shopping-lists [post]
// @Security BearerAuth
func (h *ShoppingListHandler) CreateList(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	list, err := h.Service.CreateList(r.Context(), userID, inventoryID, req.Name)
	if err != nil {
		if err.Error() == "unauthorized" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(list)
}

// ListLists lists all shopping lists.
// @Summary List shopping lists
// @Description Retrieve a list of shopping lists for an inventory.
// @Tags Shopping List
// @Produce json
// @Param id path string true "Inventory ID"
// @Success 200 {array} models.ShoppingList
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 500 {string} string "internal server error"
// @Router /inventories/{id}/shopping-lists [get]
// @Security BearerAuth
func (h *ShoppingListHandler) ListLists(w http.ResponseWriter, r *http.Request) {
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

	lists, err := h.Service.ListLists(r.Context(), userID, inventoryID)
	if err != nil {
		if err.Error() == "unauthorized" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(lists)
}

// GetList retrieves a shopping list by ID.
// @Summary Get shopping list
// @Description Get details of a specific shopping list.
// @Tags Shopping List
// @Produce json
// @Param id path string true "Shopping List ID"
// @Success 200 {object} models.ShoppingList
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "list not found"
// @Failure 500 {string} string "internal server error"
// @Router /shopping-lists/{id} [get]
// @Security BearerAuth
func (h *ShoppingListHandler) GetList(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	listID := r.PathValue("id")
	if listID == "" {
		http.Error(w, "list id required", http.StatusBadRequest)
		return
	}

	list, err := h.Service.GetList(r.Context(), userID, listID)
	if err != nil {
		if err.Error() == "unauthorized" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch items too or dedicated endpoint?
	// Usually GET list should include items or allow fetching them.
	// For API RESTfulness, GET /shopping-lists/:id/items is typical for items.
	// But getting list metadata is fine here.
	json.NewEncoder(w).Encode(list)
}

// UpdateList updates a shopping list.
// @Summary Update shopping list
// @Description Update details of an existing shopping list.
// @Tags Shopping List
// @Accept json
// @Produce json
// @Param id path string true "Shopping List ID"
// @Param request body CreateListRequest true "Update List Request"
// @Success 200 {object} models.ShoppingList
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "list not found"
// @Failure 500 {string} string "internal server error"
// @Router /shopping-lists/{id} [put]
// @Security BearerAuth
func (h *ShoppingListHandler) UpdateList(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	listID := r.PathValue("id")
	if listID == "" {
		http.Error(w, "list id required", http.StatusBadRequest)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	list, err := h.Service.UpdateList(r.Context(), userID, listID, req.Name)
	if err != nil {
		if err.Error() == "unauthorized" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(list)
}

// DeleteList deletes a shopping list.
// @Summary Delete shopping list
// @Description Delete a shopping list.
// @Tags Shopping List
// @Param id path string true "Shopping List ID"
// @Success 204
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "list not found"
// @Failure 500 {string} string "internal server error"
// @Router /shopping-lists/{id} [delete]
// @Security BearerAuth
func (h *ShoppingListHandler) DeleteList(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	listID := r.PathValue("id")
	if listID == "" {
		http.Error(w, "list id required", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteList(r.Context(), userID, listID); err != nil {
		if err.Error() == "unauthorized" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ListItems lists items in a shopping list.
// @Summary List shopping list items
// @Description Retrieve a list of items in the specified shopping list.
// @Tags Shopping List
// @Produce json
// @Param id path string true "Shopping List ID"
// @Success 200 {array} models.ShoppingListItem
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "list not found"
// @Failure 500 {string} string "internal server error"
// @Router /shopping-lists/{id}/items [get]
// @Security BearerAuth
func (h *ShoppingListHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	listID := r.PathValue("id")
	if listID == "" {
		http.Error(w, "list id required", http.StatusBadRequest)
		return
	}

	items, err := h.Service.ListItems(r.Context(), userID, listID)
	if err != nil {
		if err.Error() == "unauthorized" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(items)
}

// AddItem adds an item to a shopping list.
// @Summary Add item to shopping list
// @Description Add a new item to the shopping list.
// @Tags Shopping List
// @Accept json
// @Produce json
// @Param id path string true "Shopping List ID"
// @Param request body AddItemRequest true "Add Item Request"
// @Success 201 {object} models.ShoppingListItem
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "list not found"
// @Failure 500 {string} string "internal server error"
// @Router /shopping-lists/{id}/items [post]
// @Security BearerAuth
func (h *ShoppingListHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	listID := r.PathValue("id")
	if listID == "" {
		http.Error(w, "list id required", http.StatusBadRequest)
		return
	}

	var req struct {
		TargetType        string   `json:"target_type"`
		TargetID          string   `json:"target_id"`
		PreferredOutletID *string  `json:"preferred_outlet_id"`
		Notes             *string  `json:"notes"`
		Quantity          *float64 `json:"quantity"`
		Unit              *string  `json:"unit"`
		ManualOrder       *int     `json:"manual_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	item := &models.ShoppingListItem{
		TargetType:        req.TargetType,
		TargetID:          req.TargetID,
		PreferredOutletID: req.PreferredOutletID,
		Notes:             req.Notes,
		Quantity:          req.Quantity,
		Unit:              req.Unit,
		ManualOrder:       req.ManualOrder,
	}

	createdItem, err := h.Service.AddItem(r.Context(), userID, listID, item)
	if err != nil {
		if err.Error() == "unauthorized" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdItem)
}

// UpdateItem updates an item in a shopping list.
// @Summary Update shopping list item
// @Description Update details of an existing shopping list item.
// @Tags Shopping List
// @Accept json
// @Produce json
// @Param itemId path string true "Item ID"
// @Param request body UpdateItemRequest true "Update Item Request"
// @Success 200 {object} models.ShoppingListItem
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "item not found"
// @Failure 500 {string} string "internal server error"
// @Router /shopping-list-items/{itemId} [put]
// @Security BearerAuth
func (h *ShoppingListHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	itemID := r.PathValue("itemId")
	if itemID == "" {
		http.Error(w, "item id required", http.StatusBadRequest)
		return
	}

	var req struct {
		Notes             *string  `json:"notes"`
		PreferredOutletID *string  `json:"preferred_outlet_id"`
		Quantity          *float64 `json:"quantity"`
		Unit              *string  `json:"unit"`
		ManualOrder       *int     `json:"manual_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updatedItem, err := h.Service.UpdateItem(r.Context(), userID, itemID, req.Notes, req.PreferredOutletID, req.Quantity, req.Unit, req.ManualOrder)
	if err != nil {
		if err.Error() == "unauthorized" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedItem)
}

// DeleteItem deletes an item from a shopping list.
// @Summary Delete shopping list item
// @Description Delete an item from the shopping list.
// @Tags Shopping List
// @Param itemId path string true "Item ID"
// @Success 204
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "item not found"
// @Failure 500 {string} string "internal server error"
// @Router /shopping-list-items/{itemId} [delete]
// @Security BearerAuth
func (h *ShoppingListHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	itemID := r.PathValue("itemId")
	if itemID == "" {
		http.Error(w, "item id required", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteItem(r.Context(), userID, itemID); err != nil {
		if err.Error() == "unauthorized" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateItemsOrder updates the manual order of multiple items in a shopping list.
// @Summary Update items order
// @Description Update the manual order of multiple items in a shopping list.
// @Tags Shopping List
// @Accept json
// @Produce json
// @Param id path string true "Shopping List ID"
// @Param request body []models.UpdateOrderItem true "Update Items Order Request"
// @Success 204
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "list not found"
// @Failure 500 {string} string "internal server error"
// @Router /shopping-lists/{id}/items/order [put]
// @Security BearerAuth
func (h *ShoppingListHandler) UpdateItemsOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	listID := r.PathValue("id")
	if listID == "" {
		http.Error(w, "list id required", http.StatusBadRequest)
		return
	}

	var req []models.UpdateOrderItem
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateItemsOrder(r.Context(), userID, listID, req); err != nil {
		if err.Error() == "unauthorized" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
