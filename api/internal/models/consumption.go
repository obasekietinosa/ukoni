package models

import (
	"context"
	"database/sql"
	"time"
	"ukoni/internal/database"
)

type ConsumptionEvent struct {
	ID                 string     `json:"id"`
	InventoryID        string     `json:"inventory_id"`
	CanonicalProductID *string    `json:"canonical_product_id,omitempty"`
	CreatedByUserID    *string    `json:"created_by_user_id,omitempty"`
	Quantity           *float64   `json:"quantity,omitempty"`
	Unit               *string    `json:"unit,omitempty"`
	Note               *string    `json:"note,omitempty"`
	Source             string     `json:"source"`
	ConsumedAt         time.Time  `json:"consumed_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}

type ConsumptionModel struct {
	DB *sql.DB
}

func (m *ConsumptionModel) Create(ctx context.Context, dbtx database.DBTX, event *ConsumptionEvent) error {
	query := `
		INSERT INTO consumption_events (
			inventory_id, canonical_product_id, created_by_user_id,
			quantity, unit, note, source, consumed_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	return dbtx.QueryRowContext(ctx, query,
		event.InventoryID,
		event.CanonicalProductID,
		event.CreatedByUserID,
		event.Quantity,
		event.Unit,
		event.Note,
		event.Source,
		event.ConsumedAt,
	).Scan(&event.ID)
}

func (m *ConsumptionModel) List(ctx context.Context, inventoryID string, limit, offset int) ([]*ConsumptionEvent, error) {
	query := `
		SELECT id, inventory_id, canonical_product_id, created_by_user_id,
		       quantity, unit, note, source, consumed_at, deleted_at
		FROM consumption_events
		WHERE inventory_id = $1 AND deleted_at IS NULL
		ORDER BY consumed_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := m.DB.QueryContext(ctx, query, inventoryID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*ConsumptionEvent
	for rows.Next() {
		var e ConsumptionEvent
		if err := rows.Scan(
			&e.ID, &e.InventoryID, &e.CanonicalProductID, &e.CreatedByUserID,
			&e.Quantity, &e.Unit, &e.Note, &e.Source, &e.ConsumedAt, &e.DeletedAt,
		); err != nil {
			return nil, err
		}
		events = append(events, &e)
	}
	return events, rows.Err()
}
