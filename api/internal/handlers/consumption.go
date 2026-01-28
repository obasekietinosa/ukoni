package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"ukoni/internal/models"
	"ukoni/internal/services"
)

type ConsumptionHandler struct {
	Service *services.ConsumptionService
}

type createConsumptionRequest struct {
	CanonicalProductID *string  `json:"canonical_product_id"`
	Quantity           *float64 `json:"quantity"`
	Unit               *string  `json:"unit"`
	Note               *string  `json:"note"`
	Source             string   `json:"source"`
	ConsumedAt         string   `json:"consumed_at"` // ISO8601 string
}

func (h *ConsumptionHandler) CreateConsumptionEvent(w http.ResponseWriter, r *http.Request) {
	inventoryID := r.PathValue("id")
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req createConsumptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var consumedAt time.Time
	if req.ConsumedAt != "" {
		parsed, err := time.Parse(time.RFC3339, req.ConsumedAt)
		if err != nil {
			http.Error(w, "invalid consumed_at format (expected RFC3339)", http.StatusBadRequest)
			return
		}
		consumedAt = parsed
	} else {
		consumedAt = time.Now()
	}

	input := services.CreateConsumptionInput{
		InventoryID:        inventoryID,
		CanonicalProductID: req.CanonicalProductID,
		CreatedByUserID:    userID,
		Quantity:           req.Quantity,
		Unit:               req.Unit,
		Note:               req.Note,
		Source:             req.Source,
		ConsumedAt:         consumedAt,
	}

	event, err := h.Service.CreateConsumption(r.Context(), input)
	if err != nil {
		if err.Error() == "user is not a member of this inventory" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

func (h *ConsumptionHandler) ListConsumptionEvents(w http.ResponseWriter, r *http.Request) {
	inventoryID := r.PathValue("id")
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	limit := 20
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	events, err := h.Service.ListConsumptionEvents(r.Context(), inventoryID, userID, limit, offset)
	if err != nil {
		if err.Error() == "user is not a member of this inventory" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return empty array instead of null
	if events == nil {
		events = []*models.ConsumptionEvent{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
