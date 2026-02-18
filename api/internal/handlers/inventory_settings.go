package handlers

import (
	"encoding/json"
	"net/http"
	"ukoni/internal/models"
	"ukoni/internal/services"
)

type InventorySettingsHandler struct {
	Service *services.InventorySettingsService
}

// GetSettings godoc
// @Summary Get inventory settings
// @Description Get inventory settings
// @Tags Settings
// @Accept json
// @Produce json
// @Param id path string true "Inventory ID"
// @Success 200 {object} models.InventorySettings
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /inventories/{id}/settings [get]
// @Security BearerAuth
func (h *InventorySettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	inventoryID := r.PathValue("id")
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	settings, err := h.Service.GetSettings(r.Context(), userID, inventoryID)
	if err != nil {
		if err.Error() == "unauthorized" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

// UpdateSettings godoc
// @Summary Update inventory settings
// @Description Update inventory settings
// @Tags Settings
// @Accept json
// @Produce json
// @Param id path string true "Inventory ID"
// @Param settings body models.InventorySettings true "Settings"
// @Success 200 {object} models.InventorySettings
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /inventories/{id}/settings [put]
// @Security BearerAuth
func (h *InventorySettingsHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	inventoryID := r.PathValue("id")
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var settings models.InventorySettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated, err := h.Service.UpdateSettings(r.Context(), userID, inventoryID, &settings)
	if err != nil {
		if err.Error() == "unauthorized" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}
