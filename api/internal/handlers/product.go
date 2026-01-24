package handlers

import (
	"encoding/json"
	"net/http"
	"ukoni/internal/services"
)

type ProductHandler struct {
	Service *services.ProductService
}

// Categories

func (h *ProductHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name             string  `json:"name"`
		ParentCategoryID *string `json:"parent_category_id,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	category, err := h.Service.CreateCategory(r.Context(), userID, req.Name, req.ParentCategoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func (h *ProductHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.Service.ListCategories(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(categories)
}

// Canonical Products

func (h *ProductHandler) CreateCanonicalProduct(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name        string  `json:"name"`
		Description *string `json:"description,omitempty"`
		CategoryID  *string `json:"category_id,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	cp, err := h.Service.CreateCanonicalProduct(r.Context(), userID, req.Name, req.Description, req.CategoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cp)
}

func (h *ProductHandler) ListCanonicalProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.Service.ListCanonicalProducts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(products)
}

// Products

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		CanonicalProductID *string `json:"canonical_product_id,omitempty"`
		Brand              *string `json:"brand,omitempty"`
		Name               string  `json:"name"`
		Description        *string `json:"description,omitempty"`
		CategoryID         *string `json:"category_id,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	p, err := h.Service.CreateProduct(r.Context(), userID, req.CanonicalProductID, req.Brand, req.Name, req.Description, req.CategoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.Service.ListProducts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(products)
}

// Product Variants

func (h *ProductHandler) CreateProductVariant(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	productID := r.PathValue("id")
	if productID == "" {
		http.Error(w, "product id required", http.StatusBadRequest)
		return
	}

	var req struct {
		VariantName string  `json:"variant_name"`
		SKU         *string `json:"sku,omitempty"`
		Unit        *string `json:"unit,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	pv, err := h.Service.CreateProductVariant(r.Context(), userID, productID, req.VariantName, req.SKU, req.Unit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pv)
}

func (h *ProductHandler) ListProductVariants(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("id")
	if productID == "" {
		http.Error(w, "product id required", http.StatusBadRequest)
		return
	}

	variants, err := h.Service.ListProductVariants(r.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(variants)
}
