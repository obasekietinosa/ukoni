package handlers

import (
	"encoding/json"
	"net/http"
	"ukoni/internal/services"
)

type MembershipHandler struct {
	Service *services.MembershipService
}

type InviteUserRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type AcceptInviteRequest struct {
	Token string `json:"token"`
}

type UpdateMemberRequest struct {
	Role string `json:"role"`
}

// InviteUser handles creating a new invitation
// @Summary Invite user to inventory
// @Description Send an invitation email to a user to join the specified inventory with a role.
// @Tags Membership
// @Accept json
// @Produce json
// @Param id path string true "Inventory ID"
// @Param request body InviteUserRequest true "Invitation Request"
// @Success 201 {object} models.Invitation
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 500 {string} string "internal server error"
// @Router /inventories/{id}/invitations [post]
// @Security BearerAuth
func (h *MembershipHandler) InviteUser(w http.ResponseWriter, r *http.Request) {
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

	var req InviteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Default role if not provided
	if req.Role == "" {
		req.Role = "viewer"
	}

	invitation, err := h.Service.InviteUser(userID, inventoryID, req.Email, req.Role)
	if err != nil {
		if err == services.ErrUnauthorized {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(invitation)
}

// AcceptInvite handles accepting an invitation
// @Summary Accept invitation
// @Description Accept an invitation using the token.
// @Tags Membership
// @Accept json
// @Produce json
// @Param id path string true "Invitation ID"
// @Param request body AcceptInviteRequest true "Accept Invitation Request"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /invitations/{id}/accept [post]
// @Security BearerAuth
func (h *MembershipHandler) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		// Acceptance generally requires the user to be logged in effectively linking the invite to their account
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	inviteID := r.PathValue("id")
	if inviteID == "" {
		http.Error(w, "invitation id required", http.StatusBadRequest)
		return
	}

	var req AcceptInviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Token == "" {
		http.Error(w, "token required", http.StatusBadRequest)
		return
	}

	if err := h.Service.AcceptInvitation(userID, inviteID, req.Token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"accepted"}`))
}

// ListMembers handles listing all members of an inventory
// @Summary List inventory members
// @Description Retrieve a list of all members in the specified inventory.
// @Tags Membership
// @Produce json
// @Param id path string true "Inventory ID"
// @Success 200 {array} models.InventoryMembership
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 500 {string} string "internal server error"
// @Router /inventories/{id}/members [get]
// @Security BearerAuth
func (h *MembershipHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
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

	members, err := h.Service.ListMembers(userID, inventoryID)
	if err != nil {
		if err == services.ErrUnauthorized {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(members)
}

// RemoveMember handles removing a member from an inventory
// @Summary Remove member
// @Description Remove a user from the inventory.
// @Tags Membership
// @Param id path string true "Inventory ID"
// @Param userId path string true "User ID"
// @Success 204
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 500 {string} string "internal server error"
// @Router /inventories/{id}/members/{userId} [delete]
// @Security BearerAuth
func (h *MembershipHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	inventoryID := r.PathValue("id")
	targetUserID := r.PathValue("userId")
	if inventoryID == "" || targetUserID == "" {
		http.Error(w, "inventory id and user id required", http.StatusBadRequest)
		return
	}

	if err := h.Service.RemoveMember(userID, inventoryID, targetUserID); err != nil {
		if err == services.ErrUnauthorized {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateMember handles updating a member's role
// @Summary Update member role
// @Description Update the role of a member in the inventory.
// @Tags Membership
// @Accept json
// @Param id path string true "Inventory ID"
// @Param userId path string true "User ID"
// @Param request body UpdateMemberRequest true "Update Member Request"
// @Success 200
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 500 {string} string "internal server error"
// @Router /inventories/{id}/members/{userId} [put]
// @Security BearerAuth
func (h *MembershipHandler) UpdateMember(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	inventoryID := r.PathValue("id")
	targetUserID := r.PathValue("userId")
	if inventoryID == "" || targetUserID == "" {
		http.Error(w, "inventory id and user id required", http.StatusBadRequest)
		return
	}

	var req UpdateMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Role == "" {
		http.Error(w, "role required", http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateMemberRole(userID, inventoryID, targetUserID, req.Role); err != nil {
		if err == services.ErrUnauthorized {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
