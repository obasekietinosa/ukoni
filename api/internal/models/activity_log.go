package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
	"ukoni/internal/database"
)

type ActivityLog struct {
	ID          string                 `json:"id"`
	InventoryID *string                `json:"inventory_id,omitempty"`
	UserID      *string                `json:"user_id,omitempty"`
	Action      string                 `json:"action"`
	EntityType  string                 `json:"entity_type,omitempty"`
	EntityID    *string                `json:"entity_id,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}

type ActivityLogModel struct {
	DB *sql.DB
}

func (m *ActivityLogModel) Create(ctx context.Context, dbtx database.DBTX, logEntry *ActivityLog) error {
	query := `
		INSERT INTO activity_logs (inventory_id, user_id, action, entity_type, entity_id, metadata)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	metadataJSON, err := json.Marshal(logEntry.Metadata)
	if err != nil {
		return err
	}

	return dbtx.QueryRowContext(ctx, query,
		logEntry.InventoryID,
		logEntry.UserID,
		logEntry.Action,
		logEntry.EntityType,
		logEntry.EntityID,
		metadataJSON,
	).Scan(&logEntry.ID, &logEntry.CreatedAt)
}
