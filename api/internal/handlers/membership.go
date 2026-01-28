package handlers

import (
	"encoding/json"
	"net/http"
	"ukoni/internal/services"
)

type MembershipHandler struct {
	Service *services.MembershipService
}

// InviteUser handles creating a new invitation
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

	var req struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
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

	var req struct {
		Token string `json:"token"`
	}
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
