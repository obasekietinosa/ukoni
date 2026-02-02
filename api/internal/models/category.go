package models

import (
	"context"
	"database/sql"
	"ukoni/internal/database"
)

type Category struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	ParentCategoryID *string `json:"parent_category_id,omitempty"`
}

type CategoryModel struct {
	DB *sql.DB
}

func (m *CategoryModel) Create(ctx context.Context, dbtx database.DBTX, category *Category) error {
	query := `
		INSERT INTO product_categories (name, parent_category_id)
		VALUES ($1, $2)
		RETURNING id
	`
	return dbtx.QueryRowContext(ctx, query,
		category.Name,
		category.ParentCategoryID,
	).Scan(&category.ID)
}

func (m *CategoryModel) List(ctx context.Context) ([]*Category, error) {
	query := `
		SELECT id, name, parent_category_id
		FROM product_categories
		WHERE deleted_at IS NULL
		ORDER BY name ASC
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []*Category{}
	for rows.Next() {
		var c Category
		if err := rows.Scan(
			&c.ID, &c.Name, &c.ParentCategoryID,
		); err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}
	return categories, rows.Err()
}
