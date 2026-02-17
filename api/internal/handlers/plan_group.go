package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"ukoni/internal/services"
)

type PlanGroupHandler struct {
	Service           *services.PlanGroupService
	MembershipService *services.MembershipService
}

type CreatePlanGroupRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type AddPlanToGroupRequest struct {
	PlanID string `json:"plan_id"`
}

type LinkShoppingListToGroupRequest struct {
	ShoppingListID string `json:"shopping_list_id"`
}

type CreateShoppingListFromGroupRequest struct {
	Name string `json:"name"`
}

// CreateGroup creates a new plan group.
// @Summary Create plan group
// @Description Create a new plan group.
// @Tags PlanGroup
// @Accept json
// @Produce json
// @Param id path string true "Inventory ID"
// @Param request body CreatePlanGroupRequest true "Create Group Request"
// @Success 201 {object} models.PlanGroup
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 500 {string} string "internal server error"
// @Router /inventories/{id}/plan-groups [post]
// @Security BearerAuth
func (h *PlanGroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
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

	if _, err := h.MembershipService.MembershipModel.GetMembership(inventoryID, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "forbidden", http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var req CreatePlanGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	group, err := h.Service.CreateGroup(r.Context(), inventoryID, req.Title, req.Description)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(group)
}

// ListGroups lists plan groups.
// @Summary List plan groups
// @Description Retrieve a list of plan groups for an inventory.
// @Tags PlanGroup
// @Produce json
// @Param id path string true "Inventory ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} models.PlanGroup
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 500 {string} string "internal server error"
// @Router /inventories/{id}/plan-groups [get]
// @Security BearerAuth
func (h *PlanGroupHandler) ListGroups(w http.ResponseWriter, r *http.Request) {
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

	if _, err := h.MembershipService.MembershipModel.GetMembership(inventoryID, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "forbidden", http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	query := r.URL.Query()
	limit := 10
	offset := 0
	if l, err := strconv.Atoi(query.Get("limit")); err == nil {
		limit = l
	}
	if o, err := strconv.Atoi(query.Get("offset")); err == nil {
		offset = o
	}

	groups, err := h.Service.ListGroups(r.Context(), inventoryID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(groups)
}

// GetGroup gets a plan group.
// @Summary Get plan group
// @Description Get a plan group by ID.
// @Tags PlanGroup
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} services.PlanGroupWithDetails
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "group not found"
// @Failure 500 {string} string "internal server error"
// @Router /plan-groups/{id} [get]
// @Security BearerAuth
func (h *PlanGroupHandler) GetGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "group id required", http.StatusBadRequest)
		return
	}

	group, err := h.Service.GetGroup(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if group == nil {
		http.Error(w, "group not found", http.StatusNotFound)
		return
	}

	if _, err := h.MembershipService.MembershipModel.GetMembership(group.InventoryID, userID); err != nil {
		http.Error(w, "group not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(group)
}

// UpdateGroup updates a plan group.
// @Summary Update plan group
// @Description Update a plan group.
// @Tags PlanGroup
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body CreatePlanGroupRequest true "Update Group Request"
// @Success 200 {object} models.PlanGroup
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "group not found"
// @Failure 500 {string} string "internal server error"
// @Router /plan-groups/{id} [put]
// @Security BearerAuth
func (h *PlanGroupHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "group id required", http.StatusBadRequest)
		return
	}

	// Permission check needs fetching the group first, but UpdateGroup in service does checks.
	// But to return 404/403 correctly we should check membership.
	// We can fetch summary first. Or rely on Service error handling if service checks membership (it doesn't usually).
	// Let's fetch group first.
	group, err := h.Service.GetGroup(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if group == nil {
		http.Error(w, "group not found", http.StatusNotFound)
		return
	}
	if _, err := h.MembershipService.MembershipModel.GetMembership(group.InventoryID, userID); err != nil {
		http.Error(w, "group not found", http.StatusNotFound)
		return
	}

	var req CreatePlanGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updated, err := h.Service.UpdateGroup(r.Context(), id, req.Title, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updated)
}

// DeleteGroup deletes a plan group.
// @Summary Delete plan group
// @Description Delete a plan group.
// @Tags PlanGroup
// @Param id path string true "Group ID"
// @Success 204
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "group not found"
// @Failure 500 {string} string "internal server error"
// @Router /plan-groups/{id} [delete]
// @Security BearerAuth
func (h *PlanGroupHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "group id required", http.StatusBadRequest)
		return
	}

	group, err := h.Service.GetGroup(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if group == nil {
		http.Error(w, "group not found", http.StatusNotFound)
		return
	}
	if _, err := h.MembershipService.MembershipModel.GetMembership(group.InventoryID, userID); err != nil {
		http.Error(w, "group not found", http.StatusNotFound)
		return
	}

	err = h.Service.DeleteGroup(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddPlanToGroup adds a plan to a group.
// @Summary Add plan to group
// @Description Add a plan to a plan group.
// @Tags PlanGroup
// @Accept json
// @Param id path string true "Group ID"
// @Param request body AddPlanToGroupRequest true "Add Plan Request"
// @Success 201
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "group not found"
// @Failure 500 {string} string "internal server error"
// @Router /plan-groups/{id}/plans [post]
// @Security BearerAuth
func (h *PlanGroupHandler) AddPlanToGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "group id required", http.StatusBadRequest)
		return
	}

	group, err := h.Service.GetGroup(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if group == nil {
		http.Error(w, "group not found", http.StatusNotFound)
		return
	}
	if _, err := h.MembershipService.MembershipModel.GetMembership(group.InventoryID, userID); err != nil {
		http.Error(w, "group not found", http.StatusNotFound)
		return
	}

	var req AddPlanToGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.Service.AddPlanToGroup(r.Context(), id, req.PlanID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// RemovePlanFromGroup removes a plan from a group.
// @Summary Remove plan from group
// @Description Remove a plan from a plan group.
// @Tags PlanGroup
// @Param id path string true "Group ID"
// @Param planId path string true "Plan ID"
// @Success 204
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "group not found"
// @Failure 500 {string} string "internal server error"
// @Router /plan-groups/{id}/plans/{planId} [delete]
// @Security BearerAuth
func (h *PlanGroupHandler) RemovePlanFromGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "group id required", http.StatusBadRequest)
		return
	}
	planID := r.PathValue("planId")
	if planID == "" {
		http.Error(w, "plan id required", http.StatusBadRequest)
		return
	}

	group, err := h.Service.GetGroup(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if group == nil {
		http.Error(w, "group not found", http.StatusNotFound)
		return
	}
	if _, err := h.MembershipService.MembershipModel.GetMembership(group.InventoryID, userID); err != nil {
		http.Error(w, "group not found", http.StatusNotFound)
		return
	}

	err = h.Service.RemovePlanFromGroup(r.Context(), id, planID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CreateShoppingListFromGroup creates a shopping list from a group.
// @Summary Create shopping list from group
// @Description Create a shopping list from a plan group.
// @Tags PlanGroup
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body CreateShoppingListFromGroupRequest false "Request"
// @Success 201 {object} models.ShoppingList
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "group not found"
// @Failure 500 {string} string "internal server error"
// @Router /plan-groups/{id}/shopping-list [post]
// @Security BearerAuth
func (h *PlanGroupHandler) CreateShoppingListFromGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "group id required", http.StatusBadRequest)
		return
	}

	group, err := h.Service.GetGroup(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if group == nil {
		http.Error(w, "group not found", http.StatusNotFound)
		return
	}
	if _, err := h.MembershipService.MembershipModel.GetMembership(group.InventoryID, userID); err != nil {
		http.Error(w, "group not found", http.StatusNotFound)
		return
	}

	var req CreateShoppingListFromGroupRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	list, err := h.Service.CreateShoppingListFromGroup(r.Context(), id, req.Name, userID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(list)
}

// LinkShoppingListToGroup links a shopping list to a group.
// @Summary Link shopping list to group
// @Description Link a shopping list to a plan group.
// @Tags PlanGroup
// @Accept json
// @Param id path string true "Group ID"
// @Param request body LinkShoppingListToGroupRequest true "Link Request"
// @Success 201
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "group not found"
// @Failure 500 {string} string "internal server error"
// @Router /plan-groups/{id}/shopping-lists [post]
// @Security BearerAuth
func (h *PlanGroupHandler) LinkShoppingListToGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "group id required", http.StatusBadRequest)
		return
	}

	group, err := h.Service.GetGroup(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if group == nil {
		http.Error(w, "group not found", http.StatusNotFound)
		return
	}
	if _, err := h.MembershipService.MembershipModel.GetMembership(group.InventoryID, userID); err != nil {
		http.Error(w, "group not found", http.StatusNotFound)
		return
	}

	var req LinkShoppingListToGroupRequest
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// UnlinkShoppingListFromGroup unlinks a shopping list from a group.
// @Summary Unlink shopping list from group
// @Description Unlink a shopping list from a plan group.
// @Tags PlanGroup
// @Param id path string true "Group ID"
// @Param listId path string true "Shopping List ID"
// @Success 204
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "group not found"
// @Failure 500 {string} string "internal server error"
// @Router /plan-groups/{id}/shopping-lists/{listId} [delete]
// @Security BearerAuth
func (h *PlanGroupHandler) UnlinkShoppingListFromGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "group id required", http.StatusBadRequest)
		return
	}
	listID := r.PathValue("listId")
	if listID == "" {
		http.Error(w, "list id required", http.StatusBadRequest)
		return
	}

	group, err := h.Service.GetGroup(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if group == nil {
		http.Error(w, "group not found", http.StatusNotFound)
		return
	}
	if _, err := h.MembershipService.MembershipModel.GetMembership(group.InventoryID, userID); err != nil {
		http.Error(w, "group not found", http.StatusNotFound)
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
