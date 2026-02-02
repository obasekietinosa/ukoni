package models

import (
	"context"
	"database/sql"
	"ukoni/internal/database"
)

type Category struct {
	ID               string  `json:"id"`
	InventoryID      string  `json:"inventory_id"`
	Name             string  `json:"name"`
	ParentCategoryID *string `json:"parent_category_id,omitempty"`
}

type CategoryModel struct {
	DB *sql.DB
}

func (m *CategoryModel) Create(ctx context.Context, dbtx database.DBTX, category *Category) error {
	query := `
		INSERT INTO product_categories (inventory_id, name, parent_category_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	return dbtx.QueryRowContext(ctx, query,
		category.InventoryID,
		category.Name,
		category.ParentCategoryID,
	).Scan(&category.ID)
}

func (m *CategoryModel) List(ctx context.Context, inventoryID string) ([]*Category, error) {
	query := `
		SELECT id, inventory_id, name, parent_category_id
		FROM product_categories
		WHERE inventory_id = $1 AND deleted_at IS NULL
		ORDER BY name ASC
	`
	rows, err := m.DB.QueryContext(ctx, query, inventoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []*Category{}
	for rows.Next() {
		var c Category
		if err := rows.Scan(
			&c.ID, &c.InventoryID, &c.Name, &c.ParentCategoryID,
		); err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}
	return categories, rows.Err()
}

func (m *CategoryModel) GetByID(ctx context.Context, id string) (*Category, error) {
	query := `
		SELECT id, inventory_id, name, parent_category_id
		FROM product_categories
		WHERE id = $1 AND deleted_at IS NULL
	`
	var c Category
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.InventoryID, &c.Name, &c.ParentCategoryID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}
