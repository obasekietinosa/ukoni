package models

import (
	"context"
	"database/sql"
	"time"
	"ukoni/internal/database"
)

type InventorySettings struct {
	InventoryID string    `json:"inventory_id"`
	LLMProvider *string   `json:"llm_provider,omitempty"`
	LLMAPIKey   *string   `json:"llm_api_key,omitempty"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type InventorySettingsModel struct {
	DB *sql.DB
}

func (m *InventorySettingsModel) Get(ctx context.Context, inventoryID string) (*InventorySettings, error) {
	query := `
		SELECT inventory_id, llm_provider, llm_api_key, updated_at
		FROM inventory_settings
		WHERE inventory_id = $1
	`
	var s InventorySettings
	err := m.DB.QueryRowContext(ctx, query, inventoryID).Scan(
		&s.InventoryID, &s.LLMProvider, &s.LLMAPIKey, &s.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (m *InventorySettingsModel) Upsert(ctx context.Context, dbtx database.DBTX, settings *InventorySettings) error {
	query := `
		INSERT INTO inventory_settings (inventory_id, llm_provider, llm_api_key, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (inventory_id) DO UPDATE
		SET llm_provider = EXCLUDED.llm_provider,
		    llm_api_key = EXCLUDED.llm_api_key,
		    updated_at = NOW()
		RETURNING updated_at
	`
	return dbtx.QueryRowContext(ctx, query,
		settings.InventoryID,
		settings.LLMProvider,
		settings.LLMAPIKey,
	).Scan(&settings.UpdatedAt)
}
