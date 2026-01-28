package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"ukoni/internal/database"
)

type CanonicalProduct struct {
	ID          string     `json:"id"`
	InventoryID string     `json:"inventory_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	CategoryID  *string    `json:"category_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type CanonicalProductModel struct {
	DB *sql.DB
}

func (m *CanonicalProductModel) Create(ctx context.Context, dbtx database.DBTX, product *CanonicalProduct) error {
	query := `
		INSERT INTO canonical_products (inventory_id, name, description, category_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return dbtx.QueryRowContext(ctx, query,
		product.InventoryID,
		product.Name,
		product.Description,
		product.CategoryID,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
}

func (m *CanonicalProductModel) GetByID(ctx context.Context, id string) (*CanonicalProduct, error) {
	query := `
		SELECT id, inventory_id, name, description, category_id, created_at, updated_at, deleted_at
		FROM canonical_products
		WHERE id = $1 AND deleted_at IS NULL
	`
	var p CanonicalProduct
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.InventoryID, &p.Name, &p.Description, &p.CategoryID, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (m *CanonicalProductModel) List(ctx context.Context, inventoryID string, limit, offset int, search string) ([]*CanonicalProduct, error) {
	query := `
		SELECT id, inventory_id, name, description, category_id, created_at, updated_at, deleted_at
		FROM canonical_products
		WHERE inventory_id = $1 AND deleted_at IS NULL
	`
	args := []interface{}{inventoryID}
	argCount := 2

	if search != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argCount)
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

	products := []*CanonicalProduct{}
	for rows.Next() {
		var p CanonicalProduct
		if err := rows.Scan(
			&p.ID, &p.InventoryID, &p.Name, &p.Description, &p.CategoryID, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, rows.Err()
}

func (m *CanonicalProductModel) Update(ctx context.Context, dbtx database.DBTX, product *CanonicalProduct) error {
	query := `
		UPDATE canonical_products
		SET name = $1, description = $2, category_id = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND deleted_at IS NULL
		RETURNING updated_at
	`
	err := dbtx.QueryRowContext(ctx, query,
		product.Name,
		product.Description,
		product.CategoryID,
		product.ID,
	).Scan(&product.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return err
	}
	return nil
}

func (m *CanonicalProductModel) Delete(ctx context.Context, dbtx database.DBTX, id string) error {
	query := `
		UPDATE canonical_products
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
