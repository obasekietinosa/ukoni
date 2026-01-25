package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"ukoni/internal/database"
)

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
	Size        *float64   `json:"size,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type ProductModel struct {
	DB *sql.DB
}

func (m *ProductModel) Create(ctx context.Context, dbtx database.DBTX, product *Product) error {
	query := `
		INSERT INTO products (canonical_product_id, brand, name, description, category_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return dbtx.QueryRowContext(ctx, query,
		product.CanonicalProductID,
		product.Brand,
		product.Name,
		product.Description,
		product.CategoryID,
	).Scan(&product.ID, &product.CreatedAt)
}

func (m *ProductModel) GetByID(ctx context.Context, id string) (*Product, error) {
	query := `
		SELECT id, canonical_product_id, brand, name, description, category_id, created_at, deleted_at
		FROM products
		WHERE id = $1 AND deleted_at IS NULL
	`
	var p Product
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.CanonicalProductID, &p.Brand, &p.Name, &p.Description, &p.CategoryID, &p.CreatedAt, &p.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (m *ProductModel) List(ctx context.Context, limit, offset int, search string) ([]*Product, error) {
	query := `
		SELECT id, canonical_product_id, brand, name, description, category_id, created_at, deleted_at
		FROM products
		WHERE deleted_at IS NULL
	`
	args := []interface{}{}
	argCount := 1

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR brand ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+search+"%")
		argCount++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []*Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(
			&p.ID, &p.CanonicalProductID, &p.Brand, &p.Name, &p.Description, &p.CategoryID, &p.CreatedAt, &p.DeletedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, rows.Err()
}

func (m *ProductModel) CreateVariant(ctx context.Context, dbtx database.DBTX, variant *ProductVariant) error {
	query := `
		INSERT INTO product_variants (product_id, variant_name, sku, unit, size)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return dbtx.QueryRowContext(ctx, query,
		variant.ProductID,
		variant.VariantName,
		variant.SKU,
		variant.Unit,
		variant.Size,
	).Scan(&variant.ID)
}

func (m *ProductModel) ListVariants(ctx context.Context, productID string) ([]*ProductVariant, error) {
	query := `
		SELECT id, product_id, variant_name, sku, unit, size, deleted_at
		FROM product_variants
		WHERE product_id = $1 AND deleted_at IS NULL
		ORDER BY variant_name ASC
	`
	rows, err := m.DB.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	variants := []*ProductVariant{}
	for rows.Next() {
		var v ProductVariant
		if err := rows.Scan(
			&v.ID, &v.ProductID, &v.VariantName, &v.SKU, &v.Unit, &v.Size, &v.DeletedAt,
		); err != nil {
			return nil, err
		}
		variants = append(variants, &v)
	}
	return variants, rows.Err()
}

func (m *ProductModel) Update(ctx context.Context, dbtx database.DBTX, product *Product) error {
	query := `
		UPDATE products
		SET brand = $1, name = $2, description = $3, category_id = $4
		WHERE id = $5 AND deleted_at IS NULL
	`
	result, err := dbtx.ExecContext(ctx, query,
		product.Brand,
		product.Name,
		product.Description,
		product.CategoryID,
		product.ID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (m *ProductModel) GetVariant(ctx context.Context, id string) (*ProductVariant, error) {
	query := `
		SELECT id, product_id, variant_name, sku, unit, size, deleted_at
		FROM product_variants
		WHERE id = $1 AND deleted_at IS NULL
	`
	var v ProductVariant
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&v.ID, &v.ProductID, &v.VariantName, &v.SKU, &v.Unit, &v.Size, &v.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &v, nil
}

func (m *ProductModel) Delete(ctx context.Context, dbtx database.DBTX, id string) error {
	query := `
		UPDATE products
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL
	`
	result, err := dbtx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
