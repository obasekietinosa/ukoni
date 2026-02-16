package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"ukoni/internal/services"
)

type PlanHandler struct {
	Service           *services.PlanService
	MembershipService *services.MembershipService
}

func (h *PlanHandler) CreatePlan(w http.ResponseWriter, r *http.Request) {
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

	// Check membership
	if _, err := h.MembershipService.MembershipModel.GetMembership(inventoryID, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "forbidden", http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var req struct {
		ParentPlanID *string `json:"parent_plan_id"`
		Title        string  `json:"title"`
		Description  string  `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	plan, err := h.Service.CreatePlan(r.Context(), inventoryID, req.ParentPlanID, req.Title, req.Description)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(plan)
}

func (h *PlanHandler) ListPlans(w http.ResponseWriter, r *http.Request) {
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

	// Check membership
	if _, err := h.MembershipService.MembershipModel.GetMembership(inventoryID, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "forbidden", http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	query := r.URL.Query()
	limitStr := query.Get("limit")
	offsetStr := query.Get("offset")
	parentPlanIDStr := query.Get("parent_plan_id")

	limit := 10
	offset := 0
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	var parentPlanID *string
	if parentPlanIDStr != "" {
		parentPlanID = &parentPlanIDStr
	}

	plans, err := h.Service.ListPlans(r.Context(), inventoryID, limit, offset, parentPlanID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(plans)
}

func (h *PlanHandler) GetPlan(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "plan id required", http.StatusBadRequest)
		return
	}

	plan, err := h.Service.GetPlan(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if plan == nil {
		http.Error(w, "plan not found", http.StatusNotFound)
		return
	}

	// Check membership
	if _, err := h.MembershipService.MembershipModel.GetMembership(plan.InventoryID, userID); err != nil {
		http.Error(w, "plan not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(plan)
}

func (h *PlanHandler) UpdatePlan(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "plan id required", http.StatusBadRequest)
		return
	}

	// Check existing and permission
	existing, err := h.Service.GetPlan(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "plan not found", http.StatusNotFound)
		return
	}
	if _, err := h.MembershipService.MembershipModel.GetMembership(existing.InventoryID, userID); err != nil {
		http.Error(w, "plan not found", http.StatusNotFound)
		return
	}

	var req struct {
		ParentPlanID *string `json:"parent_plan_id"`
		Title        string  `json:"title"`
		Description  string  `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	plan, err := h.Service.UpdatePlan(r.Context(), id, req.Title, req.Description, req.ParentPlanID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			http.Error(w, "plan not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(plan)
}

func (h *PlanHandler) DeletePlan(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "plan id required", http.StatusBadRequest)
		return
	}

	// Check existing and permission
	existing, err := h.Service.GetPlan(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "plan not found", http.StatusNotFound)
		return
	}
	if _, err := h.MembershipService.MembershipModel.GetMembership(existing.InventoryID, userID); err != nil {
		http.Error(w, "plan not found", http.StatusNotFound)
		return
	}

	err = h.Service.DeletePlan(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			http.Error(w, "plan not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PlanHandler) AddPlanItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "plan id required", http.StatusBadRequest)
		return
	}

	// Check existing and permission
	existing, err := h.Service.GetPlan(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "plan not found", http.StatusNotFound)
		return
	}
	if _, err := h.MembershipService.MembershipModel.GetMembership(existing.InventoryID, userID); err != nil {
		http.Error(w, "plan not found", http.StatusNotFound)
		return
	}

	var req struct {
		TargetType string   `json:"target_type"`
		TargetID   string   `json:"target_id"`
		Quantity   *float64 `json:"quantity"`
		Unit       string   `json:"unit"`
		Note       string   `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	item, err := h.Service.AddItem(r.Context(), id, req.TargetType, req.TargetID, req.Quantity, req.Unit, req.Note)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			http.Error(w, "plan not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *PlanHandler) UpdatePlanItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "item id required", http.StatusBadRequest)
		return
	}

	// Check existing item to find plan and check permission
	item, err := h.Service.GetItem(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if item == nil {
		http.Error(w, "item not found", http.StatusNotFound)
		return
	}

	plan, err := h.Service.GetPlan(r.Context(), item.PlanID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if plan == nil {
		// Inconsistent state
		http.Error(w, "plan not found", http.StatusNotFound)
		return
	}

	if _, err := h.MembershipService.MembershipModel.GetMembership(plan.InventoryID, userID); err != nil {
		http.Error(w, "item not found", http.StatusNotFound)
		return
	}

	var req struct {
		Quantity *float64 `json:"quantity"`
		Unit     string   `json:"unit"`
		Note     string   `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updatedItem, err := h.Service.UpdateItem(r.Context(), id, req.Quantity, req.Unit, req.Note)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			http.Error(w, "item not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedItem)
}

func (h *PlanHandler) RemovePlanItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "item id required", http.StatusBadRequest)
		return
	}

	// Check existing item to find plan and check permission
	item, err := h.Service.GetItem(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if item == nil {
		http.Error(w, "item not found", http.StatusNotFound)
		return
	}

	plan, err := h.Service.GetPlan(r.Context(), item.PlanID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if plan == nil {
		http.Error(w, "plan not found", http.StatusNotFound)
		return
	}

	if _, err := h.MembershipService.MembershipModel.GetMembership(plan.InventoryID, userID); err != nil {
		http.Error(w, "item not found", http.StatusNotFound)
		return
	}

	err = h.Service.RemoveItem(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			http.Error(w, "item not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PlanHandler) LinkShoppingList(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "plan id required", http.StatusBadRequest)
		return
	}

	// Check existing plan and permission
	existing, err := h.Service.GetPlan(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "plan not found", http.StatusNotFound)
		return
	}
	if _, err := h.MembershipService.MembershipModel.GetMembership(existing.InventoryID, userID); err != nil {
		http.Error(w, "plan not found", http.StatusNotFound)
		return
	}

	var req struct {
		ShoppingListID string `json:"shopping_list_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.Service.LinkShoppingList(r.Context(), id, req.ShoppingListID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			http.Error(w, "plan or shopping list not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *PlanHandler) UnlinkShoppingList(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "plan id required", http.StatusBadRequest)
		return
	}

	listID := r.PathValue("listId")
	if listID == "" {
		http.Error(w, "list id required", http.StatusBadRequest)
		return
	}

	// Check existing plan and permission
	existing, err := h.Service.GetPlan(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "plan not found", http.StatusNotFound)
		return
	}
	if _, err := h.MembershipService.MembershipModel.GetMembership(existing.InventoryID, userID); err != nil {
		http.Error(w, "plan not found", http.StatusNotFound)
		return
	}

	err = h.Service.UnlinkShoppingList(r.Context(), id, listID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PlanHandler) CreateShoppingListFromGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "plan id required", http.StatusBadRequest)
		return
	}

	// Check existing plan and permission
	existing, err := h.Service.GetPlan(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "plan not found", http.StatusNotFound)
		return
	}
	if _, err := h.MembershipService.MembershipModel.GetMembership(existing.InventoryID, userID); err != nil {
		http.Error(w, "plan not found", http.StatusNotFound)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	// Decode is optional, as name is optional
	_ = json.NewDecoder(r.Body).Decode(&req)

	list, err := h.Service.CreateShoppingListFromGroup(r.Context(), id, req.Name, userID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			http.Error(w, "plan not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(list)
}
