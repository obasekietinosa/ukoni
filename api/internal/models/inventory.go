package models

import (
	"context"
	"database/sql"
	"time"
	"ukoni/internal/database"
)

type Inventory struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	OwnerUserID string     `json:"owner_user_id"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type InventoryModel struct {
	DB *sql.DB
}

func (m *InventoryModel) Create(ctx context.Context, dbtx database.DBTX, inventory *Inventory) error {
	query := `
		INSERT INTO inventories (name, owner_user_id)
		VALUES ($1, $2)
		RETURNING id, created_at
	`
	return dbtx.QueryRowContext(ctx, query, inventory.Name, inventory.OwnerUserID).
		Scan(&inventory.ID, &inventory.CreatedAt)
}

func (m *InventoryModel) GetByID(id string) (*Inventory, error) {
	query := `
		SELECT id, name, owner_user_id, created_at, deleted_at
		FROM inventories
		WHERE id = $1 AND deleted_at IS NULL
	`
	var i Inventory
	err := m.DB.QueryRowContext(context.Background(), query, id).Scan(
		&i.ID, &i.Name, &i.OwnerUserID, &i.CreatedAt, &i.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (m *InventoryModel) ListByUserID(userID string) ([]*Inventory, error) {
	query := `
		SELECT i.id, i.name, i.owner_user_id, i.created_at, i.deleted_at
		FROM inventories i
		LEFT JOIN inventory_memberships im ON i.id = im.inventory_id
		WHERE (i.owner_user_id = $1 OR im.user_id = $1) AND i.deleted_at IS NULL
	`
	rows, err := m.DB.QueryContext(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventories []*Inventory
	for rows.Next() {
		var i Inventory
		if err := rows.Scan(&i.ID, &i.Name, &i.OwnerUserID, &i.CreatedAt, &i.DeletedAt); err != nil {
			return nil, err
		}
		inventories = append(inventories, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return inventories, nil
}

// Ensure UUID validity check helper if needed, but for now assuming valid UUID strings from higher layers or DB handles generation.
// Actually, input validation should happen in service/handler layer.
