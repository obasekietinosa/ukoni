package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"ukoni/internal/services"
)

type ProductHandler struct {
	Service *services.ProductService
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Brand       string `json:"brand"`
		Name        string `json:"name"`
		Description string `json:"description"`
		CategoryID  string `json:"category_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	product, err := h.Service.CreateProduct(r.Context(), req.Brand, req.Name, req.Description, req.CategoryID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "product id required", http.StatusBadRequest)
		return
	}

	product, err := h.Service.GetProduct(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if product == nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "product id required", http.StatusBadRequest)
		return
	}

	var req struct {
		Brand       string `json:"brand"`
		Name        string `json:"name"`
		Description string `json:"description"`
		CategoryID  string `json:"category_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	product, err := h.Service.UpdateProduct(r.Context(), id, req.Brand, req.Name, req.Description, req.CategoryID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			http.Error(w, "product not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if product == nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "product id required", http.StatusBadRequest)
		return
	}

	err := h.Service.DeleteProduct(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			http.Error(w, "product not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limitStr := query.Get("limit")
	offsetStr := query.Get("offset")
	search := query.Get("search")

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

	products, err := h.Service.ListProducts(r.Context(), limit, offset, search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) CreateVariant(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("id")
	if productID == "" {
		http.Error(w, "product id required", http.StatusBadRequest)
		return
	}

	var req struct {
		VariantName string `json:"variant_name"`
		SKU         string `json:"sku"`
		Unit        string `json:"unit"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	variant, err := h.Service.CreateVariant(r.Context(), productID, req.VariantName, req.SKU, req.Unit)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(variant)
}

func (h *ProductHandler) ListVariants(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("id")
	if productID == "" {
		http.Error(w, "product id required", http.StatusBadRequest)
		return
	}

	variants, err := h.Service.ListVariants(r.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(variants)
}
