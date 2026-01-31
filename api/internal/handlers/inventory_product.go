package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"ukoni/internal/models"
	"ukoni/internal/services"
)

type InventoryProductHandler struct {
	Service           *services.InventoryProductService
	MembershipService *services.MembershipService
}

func (h *InventoryProductHandler) ListInventoryProducts(w http.ResponseWriter, r *http.Request) {
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
	isMember, err := h.MembershipService.IsMember(r.Context(), inventoryID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isMember {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	limit := 50
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil {
			offset = v
		}
	}

	products, err := h.Service.ListInventoryProducts(r.Context(), inventoryID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if products == nil {
		products = []*models.InventoryProductDetail{}
	}

	json.NewEncoder(w).Encode(products)
}
