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
		TargetType        string  `json:"target_type"`
		TargetID          string  `json:"target_id"`
		PreferredOutletID *string `json:"preferred_outlet_id"`
		Notes             *string `json:"notes"`
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
		Notes             *string `json:"notes"`
		PreferredOutletID *string `json:"preferred_outlet_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updatedItem, err := h.Service.UpdateItem(r.Context(), userID, itemID, req.Notes, req.PreferredOutletID)
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
