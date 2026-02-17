package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"ukoni/internal/services"
)

type ProductHandler struct {
	Service           *services.ProductService
	MembershipService *services.MembershipService
}

type ProductRequest struct {
	Brand              string `json:"brand"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	CategoryID         string `json:"category_id"`
	CanonicalProductID string `json:"canonical_product_id"`
}

type VariantRequest struct {
	VariantName string   `json:"variant_name"`
	SKU         string   `json:"sku"`
	Unit        string   `json:"unit"`
	Size        *float64 `json:"size"`
}

// CreateProduct creates a new product.
// @Summary Create product
// @Description Create a new product in the inventory.
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "Inventory ID"
// @Param request body ProductRequest true "Product Request"
// @Success 201 {object} models.Product
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 500 {string} string "internal server error"
// @Router /inventories/{id}/products [post]
// @Security BearerAuth
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
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

	var req ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	product, err := h.Service.CreateProduct(r.Context(), inventoryID, req.Brand, req.Name, req.Description, req.CategoryID, req.CanonicalProductID)
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

// GetProduct retrieves a product by ID.
// @Summary Get product
// @Description Get details of a specific product.
// @Tags Product
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "product not found"
// @Failure 500 {string} string "internal server error"
// @Router /products/{id} [get]
// @Security BearerAuth
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

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

	// Check membership
	if _, err := h.MembershipService.MembershipModel.GetMembership(product.InventoryID, userID); err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}

// UpdateProduct updates a product.
// @Summary Update product
// @Description Update details of an existing product.
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body ProductRequest true "Product Request"
// @Success 200 {object} models.Product
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "product not found"
// @Failure 500 {string} string "internal server error"
// @Router /products/{id} [put]
// @Security BearerAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "product id required", http.StatusBadRequest)
		return
	}

	// Check existence and permission first
	existing, err := h.Service.GetProduct(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}
	if _, err := h.MembershipService.MembershipModel.GetMembership(existing.InventoryID, userID); err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	var req ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	product, err := h.Service.UpdateProduct(r.Context(), id, req.Brand, req.Name, req.Description, req.CategoryID, req.CanonicalProductID)
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

// DeleteProduct deletes a product.
// @Summary Delete product
// @Description Delete a product from the inventory.
// @Tags Product
// @Param id path string true "Product ID"
// @Success 204
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "product not found"
// @Failure 500 {string} string "internal server error"
// @Router /products/{id} [delete]
// @Security BearerAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "product id required", http.StatusBadRequest)
		return
	}

	// Check existence and permission first
	existing, err := h.Service.GetProduct(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}
	if _, err := h.MembershipService.MembershipModel.GetMembership(existing.InventoryID, userID); err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	err = h.Service.DeleteProduct(r.Context(), id)
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

// ListProducts lists products in an inventory.
// @Summary List products
// @Description Retrieve a list of products in the specified inventory with pagination and filtering.
// @Tags Product
// @Produce json
// @Param id path string true "Inventory ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param search query string false "Search term"
// @Param canonical_product_id query string false "Canonical Product ID"
// @Param category_id query string false "Category ID"
// @Success 200 {array} models.Product
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 500 {string} string "internal server error"
// @Router /inventories/{id}/products [get]
// @Security BearerAuth
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
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
	search := query.Get("search")
	canonicalProductID := query.Get("canonical_product_id")
	categoryID := query.Get("category_id")

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

	products, err := h.Service.ListProducts(r.Context(), inventoryID, limit, offset, search, canonicalProductID, categoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}

// CreateVariant creates a new product variant.
// @Summary Create product variant
// @Description Create a new variant for a product.
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body VariantRequest true "Variant Request"
// @Success 201 {object} models.ProductVariant
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "product not found"
// @Failure 500 {string} string "internal server error"
// @Router /products/{id}/variants [post]
// @Security BearerAuth
func (h *ProductHandler) CreateVariant(w http.ResponseWriter, r *http.Request) {
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

	// Check existence and permission first
	product, err := h.Service.GetProduct(r.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if product == nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}
	if _, err := h.MembershipService.MembershipModel.GetMembership(product.InventoryID, userID); err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	var req VariantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	variant, err := h.Service.CreateVariant(r.Context(), productID, req.VariantName, req.SKU, req.Unit, req.Size)
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

// ListVariants lists variants of a product.
// @Summary List product variants
// @Description Retrieve a list of variants for a product.
// @Tags Product
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {array} models.ProductVariant
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "product not found"
// @Failure 500 {string} string "internal server error"
// @Router /products/{id}/variants [get]
// @Security BearerAuth
func (h *ProductHandler) ListVariants(w http.ResponseWriter, r *http.Request) {
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

	// Check existence and permission first
	product, err := h.Service.GetProduct(r.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if product == nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}
	if _, err := h.MembershipService.MembershipModel.GetMembership(product.InventoryID, userID); err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	variants, err := h.Service.ListVariants(r.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(variants)
}

// GetVariant retrieves a product variant by ID.
// @Summary Get product variant
// @Description Get details of a specific product variant.
// @Tags Product
// @Produce json
// @Param id path string true "Variant ID"
// @Success 200 {object} models.ProductVariant
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "variant not found"
// @Failure 500 {string} string "internal server error"
// @Router /product-variants/{id} [get]
// @Security BearerAuth
func (h *ProductHandler) GetVariant(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "variant id required", http.StatusBadRequest)
		return
	}

	variant, err := h.Service.GetVariant(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if variant == nil {
		http.Error(w, "variant not found", http.StatusNotFound)
		return
	}

	// Check permissions via product -> inventory
	product, err := h.Service.GetProduct(r.Context(), variant.ProductID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if product == nil {
		// Should not happen if foreign keys are enforced, but handle gracefully
		http.Error(w, "variant not found", http.StatusNotFound)
		return
	}
	if _, err := h.MembershipService.MembershipModel.GetMembership(product.InventoryID, userID); err != nil {
		http.Error(w, "variant not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(variant)
}

// UpdateVariant updates a product variant.
// @Summary Update product variant
// @Description Update details of an existing product variant.
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "Variant ID"
// @Param request body VariantRequest true "Variant Request"
// @Success 200 {object} models.ProductVariant
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "variant not found"
// @Failure 500 {string} string "internal server error"
// @Router /product-variants/{id} [put]
// @Security BearerAuth
func (h *ProductHandler) UpdateVariant(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "variant id required", http.StatusBadRequest)
		return
	}

	// Check existence and permission first
	existing, err := h.Service.GetVariant(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "variant not found", http.StatusNotFound)
		return
	}

	product, err := h.Service.GetProduct(r.Context(), existing.ProductID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if product == nil {
		http.Error(w, "variant not found", http.StatusNotFound)
		return
	}
	if _, err := h.MembershipService.MembershipModel.GetMembership(product.InventoryID, userID); err != nil {
		http.Error(w, "variant not found", http.StatusNotFound)
		return
	}

	var req VariantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	variant, err := h.Service.UpdateVariant(r.Context(), id, req.VariantName, req.SKU, req.Unit, req.Size)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			http.Error(w, "variant not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(variant)
}

// DeleteVariant deletes a product variant.
// @Summary Delete product variant
// @Description Delete a product variant.
// @Tags Product
// @Param id path string true "Variant ID"
// @Success 204
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "variant not found"
// @Failure 500 {string} string "internal server error"
// @Router /product-variants/{id} [delete]
// @Security BearerAuth
func (h *ProductHandler) DeleteVariant(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "variant id required", http.StatusBadRequest)
		return
	}

	// Check existence and permission first
	existing, err := h.Service.GetVariant(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "variant not found", http.StatusNotFound)
		return
	}

	product, err := h.Service.GetProduct(r.Context(), existing.ProductID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if product == nil {
		http.Error(w, "variant not found", http.StatusNotFound)
		return
	}
	if _, err := h.MembershipService.MembershipModel.GetMembership(product.InventoryID, userID); err != nil {
		http.Error(w, "variant not found", http.StatusNotFound)
		return
	}

	err = h.Service.DeleteVariant(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			http.Error(w, "variant not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
