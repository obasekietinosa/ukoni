package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"ukoni/internal/models"
	"ukoni/internal/services"
)

type TransactionHandler struct {
	Service *services.TransactionService
}

type CreateTransactionRequest struct {
	OutletID        *string                        `json:"outlet_id,omitempty"`
	TransactionDate time.Time                      `json:"transaction_date"`
	Items           []CreateTransactionItemRequest `json:"items"`
}

type CreateTransactionItemRequest struct {
	ProductVariantID   string   `json:"product_variant_id"`
	Quantity           float64  `json:"quantity"`
	PricePerUnit       *float64 `json:"price_per_unit,omitempty"`
	ShoppingListItemID *string  `json:"shopping_list_item_id,omitempty"`
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
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

	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	input := services.CreateTransactionInput{
		InventoryID:     inventoryID,
		CreatedByUserID: userID,
		OutletID:        req.OutletID,
		TransactionDate: req.TransactionDate,
	}

	for _, itemReq := range req.Items {
		input.Items = append(input.Items, services.CreateTransactionItemInput{
			ProductVariantID:   itemReq.ProductVariantID,
			Quantity:           itemReq.Quantity,
			PricePerUnit:       itemReq.PricePerUnit,
			ShoppingListItemID: itemReq.ShoppingListItemID,
		})
	}

	transaction, err := h.Service.CreateTransaction(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

func (h *TransactionHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
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

	limit := 10
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

	transactions, err := h.Service.ListTransactions(r.Context(), inventoryID, userID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if transactions == nil {
		transactions = []*models.Transaction{}
	}

	json.NewEncoder(w).Encode(transactions)
}

func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	transactionID := r.PathValue("id")
	if transactionID == "" {
		http.Error(w, "transaction id required", http.StatusBadRequest)
		return
	}

	transaction, err := h.Service.GetTransaction(r.Context(), transactionID, userID)
	if err != nil {
		if err == services.ErrTransactionNotFound {
			http.Error(w, "transaction not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transaction)
}
