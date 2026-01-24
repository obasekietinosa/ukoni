package models

import (
	"context"
	"database/sql"
	"time"
	"ukoni/internal/database"
)

type ProductCategory struct {
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	ParentCategoryID *string    `json:"parent_category_id,omitempty"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

type CanonicalProduct struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	CategoryID  *string    `json:"category_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type Product struct {
	ID                 string     `json:"id"`
	CanonicalProductID *string    `json:"canonical_product_id,omitempty"`
	Brand              *string    `json:"brand,omitempty"`
	Name               string     `json:"name"`
	Description        *string    `json:"description,omitempty"`
	CategoryID         *string    `json:"category_id,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}

type ProductVariant struct {
	ID          string     `json:"id"`
	ProductID   string     `json:"product_id"`
	VariantName string     `json:"variant_name"`
	SKU         *string    `json:"sku,omitempty"`
	Unit        *string    `json:"unit,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type ProductModel struct {
	DB *sql.DB
}

// ProductCategory Methods

func (m *ProductModel) CreateCategory(ctx context.Context, dbtx database.DBTX, category *ProductCategory) error {
	query := `
		INSERT INTO product_categories (name, parent_category_id)
		VALUES ($1, $2)
		RETURNING id
	`
	return dbtx.QueryRowContext(ctx, query, category.Name, category.ParentCategoryID).Scan(&category.ID)
}

func (m *ProductModel) GetCategory(ctx context.Context, id string) (*ProductCategory, error) {
	query := `SELECT id, name, parent_category_id, deleted_at FROM product_categories WHERE id = $1 AND deleted_at IS NULL`
	var c ProductCategory
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.Name, &c.ParentCategoryID, &c.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (m *ProductModel) ListCategories(ctx context.Context) ([]*ProductCategory, error) {
	query := `SELECT id, name, parent_category_id, deleted_at FROM product_categories WHERE deleted_at IS NULL ORDER BY name`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []*ProductCategory
	for rows.Next() {
		var c ProductCategory
		if err := rows.Scan(&c.ID, &c.Name, &c.ParentCategoryID, &c.DeletedAt); err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}
	return categories, nil
}

// CanonicalProduct Methods

func (m *ProductModel) CreateCanonicalProduct(ctx context.Context, dbtx database.DBTX, cp *CanonicalProduct) error {
	query := `
		INSERT INTO canonical_products (name, description, category_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	return dbtx.QueryRowContext(ctx, query, cp.Name, cp.Description, cp.CategoryID).Scan(&cp.ID, &cp.CreatedAt, &cp.UpdatedAt)
}

func (m *ProductModel) GetCanonicalProduct(ctx context.Context, id string) (*CanonicalProduct, error) {
	query := `SELECT id, name, description, category_id, created_at, updated_at, deleted_at FROM canonical_products WHERE id = $1 AND deleted_at IS NULL`
	var cp CanonicalProduct
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&cp.ID, &cp.Name, &cp.Description, &cp.CategoryID, &cp.CreatedAt, &cp.UpdatedAt, &cp.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &cp, nil
}

func (m *ProductModel) ListCanonicalProducts(ctx context.Context) ([]*CanonicalProduct, error) {
	query := `SELECT id, name, description, category_id, created_at, updated_at, deleted_at FROM canonical_products WHERE deleted_at IS NULL ORDER BY name`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []*CanonicalProduct
	for rows.Next() {
		var cp CanonicalProduct
		if err := rows.Scan(&cp.ID, &cp.Name, &cp.Description, &cp.CategoryID, &cp.CreatedAt, &cp.UpdatedAt, &cp.DeletedAt); err != nil {
			return nil, err
		}
		products = append(products, &cp)
	}
	return products, nil
}

// Product Methods

func (m *ProductModel) CreateProduct(ctx context.Context, dbtx database.DBTX, p *Product) error {
	query := `
		INSERT INTO products (canonical_product_id, brand, name, description, category_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return dbtx.QueryRowContext(ctx, query, p.CanonicalProductID, p.Brand, p.Name, p.Description, p.CategoryID).Scan(&p.ID, &p.CreatedAt)
}

func (m *ProductModel) GetProduct(ctx context.Context, id string) (*Product, error) {
	query := `SELECT id, canonical_product_id, brand, name, description, category_id, created_at, deleted_at FROM products WHERE id = $1 AND deleted_at IS NULL`
	var p Product
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.CanonicalProductID, &p.Brand, &p.Name, &p.Description, &p.CategoryID, &p.CreatedAt, &p.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (m *ProductModel) ListProducts(ctx context.Context) ([]*Product, error) {
	query := `SELECT id, canonical_product_id, brand, name, description, category_id, created_at, deleted_at FROM products WHERE deleted_at IS NULL ORDER BY name`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []*Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.CanonicalProductID, &p.Brand, &p.Name, &p.Description, &p.CategoryID, &p.CreatedAt, &p.DeletedAt); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, nil
}

// ProductVariant Methods

func (m *ProductModel) CreateProductVariant(ctx context.Context, dbtx database.DBTX, pv *ProductVariant) error {
	query := `
		INSERT INTO product_variants (product_id, variant_name, sku, unit)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	return dbtx.QueryRowContext(ctx, query, pv.ProductID, pv.VariantName, pv.SKU, pv.Unit).Scan(&pv.ID)
}

func (m *ProductModel) GetProductVariant(ctx context.Context, id string) (*ProductVariant, error) {
	query := `SELECT id, product_id, variant_name, sku, unit, deleted_at FROM product_variants WHERE id = $1 AND deleted_at IS NULL`
	var pv ProductVariant
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&pv.ID, &pv.ProductID, &pv.VariantName, &pv.SKU, &pv.Unit, &pv.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &pv, nil
}

func (m *ProductModel) ListProductVariants(ctx context.Context, productID string) ([]*ProductVariant, error) {
	query := `SELECT id, product_id, variant_name, sku, unit, deleted_at FROM product_variants WHERE product_id = $1 AND deleted_at IS NULL ORDER BY variant_name`
	rows, err := m.DB.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var variants []*ProductVariant
	for rows.Next() {
		var pv ProductVariant
		if err := rows.Scan(&pv.ID, &pv.ProductID, &pv.VariantName, &pv.SKU, &pv.Unit, &pv.DeletedAt); err != nil {
			return nil, err
		}
		variants = append(variants, &pv)
	}
	return variants, nil
}
