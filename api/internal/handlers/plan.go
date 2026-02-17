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

type CreatePlanRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type PlanItemRequest struct {
	TargetType string   `json:"target_type"`
	TargetID   string   `json:"target_id"`
	Quantity   *float64 `json:"quantity"`
	Unit       string   `json:"unit"`
	Note       string   `json:"note"`
}

type UpdatePlanItemRequest struct {
	Quantity *float64 `json:"quantity"`
	Unit     string   `json:"unit"`
	Note     string   `json:"note"`
}

type LinkShoppingListRequest struct {
	ShoppingListID string `json:"shopping_list_id"`
}

type CreateShoppingListRequest struct {
	Name string `json:"name"`
}

// CreatePlan creates a new plan.
// @Summary Create plan
// @Description Create a new plan (e.g. meal plan, event).
// @Tags Plan
// @Accept json
// @Produce json
// @Param id path string true "Inventory ID"
// @Param request body CreatePlanRequest true "Create Plan Request"
// @Success 201 {object} models.Plan
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 500 {string} string "internal server error"
// @Router /inventories/{id}/plans [post]
// @Security BearerAuth
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

	var req CreatePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	plan, err := h.Service.CreatePlan(r.Context(), inventoryID, req.Title, req.Description)
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

// ListPlans lists plans.
// @Summary List plans
// @Description Retrieve a list of plans for an inventory.
// @Tags Plan
// @Produce json
// @Param id path string true "Inventory ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} models.Plan
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 500 {string} string "internal server error"
// @Router /inventories/{id}/plans [get]
// @Security BearerAuth
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

	plans, err := h.Service.ListPlans(r.Context(), inventoryID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(plans)
}

// GetPlan retrieves a plan by ID.
// @Summary Get plan
// @Description Get details of a specific plan.
// @Tags Plan
// @Produce json
// @Param id path string true "Plan ID"
// @Success 200 {object} models.Plan
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "plan not found"
// @Failure 500 {string} string "internal server error"
// @Router /plans/{id} [get]
// @Security BearerAuth
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

// UpdatePlan updates a plan.
// @Summary Update plan
// @Description Update details of an existing plan.
// @Tags Plan
// @Accept json
// @Produce json
// @Param id path string true "Plan ID"
// @Param request body CreatePlanRequest true "Update Plan Request"
// @Success 200 {object} models.Plan
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "plan not found"
// @Failure 500 {string} string "internal server error"
// @Router /plans/{id} [put]
// @Security BearerAuth
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

	var req CreatePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	plan, err := h.Service.UpdatePlan(r.Context(), id, req.Title, req.Description)
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

// DeletePlan deletes a plan.
// @Summary Delete plan
// @Description Delete a plan.
// @Tags Plan
// @Param id path string true "Plan ID"
// @Success 204
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "plan not found"
// @Failure 500 {string} string "internal server error"
// @Router /plans/{id} [delete]
// @Security BearerAuth
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

// AddPlanItem adds an item to a plan.
// @Summary Add item to plan
// @Description Add a new item to the plan.
// @Tags Plan
// @Accept json
// @Produce json
// @Param id path string true "Plan ID"
// @Param request body PlanItemRequest true "Plan Item Request"
// @Success 201 {object} models.PlanItem
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "plan not found"
// @Failure 500 {string} string "internal server error"
// @Router /plans/{id}/items [post]
// @Security BearerAuth
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

	var req PlanItemRequest
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

// UpdatePlanItem updates a plan item.
// @Summary Update plan item
// @Description Update details of an existing plan item.
// @Tags Plan
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param request body UpdatePlanItemRequest true "Update Plan Item Request"
// @Success 200 {object} models.PlanItem
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "item not found"
// @Failure 500 {string} string "internal server error"
// @Router /plan-items/{id} [put]
// @Security BearerAuth
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

	plan, err := h.Service.GetPlanSummary(r.Context(), item.PlanID)
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

	var req UpdatePlanItemRequest
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

// RemovePlanItem removes an item from a plan.
// @Summary Remove plan item
// @Description Remove an item from the plan.
// @Tags Plan
// @Param id path string true "Item ID"
// @Success 204
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "item not found"
// @Failure 500 {string} string "internal server error"
// @Router /plan-items/{id} [delete]
// @Security BearerAuth
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

	plan, err := h.Service.GetPlanSummary(r.Context(), item.PlanID)
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

// LinkShoppingList links a shopping list to a plan.
// @Summary Link shopping list
// @Description Link a shopping list to the plan.
// @Tags Plan
// @Accept json
// @Param id path string true "Plan ID"
// @Param request body LinkShoppingListRequest true "Link Shopping List Request"
// @Success 201
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "plan not found"
// @Failure 500 {string} string "internal server error"
// @Router /plans/{id}/shopping-lists [post]
// @Security BearerAuth
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
	existing, err := h.Service.GetPlanSummary(r.Context(), id)
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

	var req LinkShoppingListRequest
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

// UnlinkShoppingList unlinks a shopping list from a plan.
// @Summary Unlink shopping list
// @Description Unlink a shopping list from the plan.
// @Tags Plan
// @Param id path string true "Plan ID"
// @Param listId path string true "Shopping List ID"
// @Success 204
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "plan not found"
// @Failure 500 {string} string "internal server error"
// @Router /plans/{id}/shopping-lists/{listId} [delete]
// @Security BearerAuth
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
	existing, err := h.Service.GetPlanSummary(r.Context(), id)
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

// CreateShoppingListFromPlan creates a shopping list from a plan.
// @Summary Create shopping list from plan
// @Description Create a shopping list containing items from the plan.
// @Tags Plan
// @Accept json
// @Produce json
// @Param id path string true "Plan ID"
// @Param request body CreateShoppingListRequest false "Create Request"
// @Success 201 {object} models.ShoppingList
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "plan not found"
// @Failure 500 {string} string "internal server error"
// @Router /plans/{id}/shopping-list [post]
// @Security BearerAuth
func (h *PlanHandler) CreateShoppingListFromPlan(w http.ResponseWriter, r *http.Request) {
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
	existing, err := h.Service.GetPlanSummary(r.Context(), id)
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

	var req CreateShoppingListRequest
	// Decode is optional, as name is optional
	_ = json.NewDecoder(r.Body).Decode(&req)

	list, err := h.Service.CreateShoppingListFromPlan(r.Context(), id, req.Name, userID)
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
