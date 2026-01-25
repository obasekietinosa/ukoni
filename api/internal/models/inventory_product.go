package models

import (
	"context"
	"database/sql"
	"time"
	"ukoni/internal/database"
)

type InventoryProduct struct {
	ID               string     `json:"id"`
	InventoryID      string     `json:"inventory_id"`
	ProductVariantID string     `json:"product_variant_id"`
	Quantity         float64    `json:"quantity"`
	Unit             *string    `json:"unit,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	LastUpdated      time.Time  `json:"last_updated"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

type InventoryProductModel struct {
	DB *sql.DB
}

func (m *InventoryProductModel) Get(ctx context.Context, inventoryID, productVariantID string) (*InventoryProduct, error) {
	query := `
		SELECT id, inventory_id, product_variant_id, quantity, unit, created_at, last_updated, deleted_at
		FROM inventory_products
		WHERE inventory_id = $1 AND product_variant_id = $2 AND deleted_at IS NULL
	`
	var ip InventoryProduct
	err := m.DB.QueryRowContext(ctx, query, inventoryID, productVariantID).Scan(
		&ip.ID, &ip.InventoryID, &ip.ProductVariantID, &ip.Quantity, &ip.Unit, &ip.CreatedAt, &ip.LastUpdated, &ip.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ip, nil
}

func (m *InventoryProductModel) Upsert(ctx context.Context, dbtx database.DBTX, inventoryID, productVariantID string, quantityChange float64, unit *string) error {
	query := `
		INSERT INTO inventory_products (inventory_id, product_variant_id, quantity, unit)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (inventory_id, product_variant_id) WHERE deleted_at IS NULL
		DO UPDATE SET
			quantity = inventory_products.quantity + EXCLUDED.quantity,
			last_updated = CURRENT_TIMESTAMP,
			unit = COALESCE(EXCLUDED.unit, inventory_products.unit)
		RETURNING id
	`
	var id string
	return dbtx.QueryRowContext(ctx, query, inventoryID, productVariantID, quantityChange, unit).Scan(&id)
}

func (m *InventoryProductModel) List(ctx context.Context, inventoryID string, limit, offset int) ([]*InventoryProduct, error) {
	query := `
		SELECT id, inventory_id, product_variant_id, quantity, unit, created_at, last_updated, deleted_at
		FROM inventory_products
		WHERE inventory_id = $1 AND deleted_at IS NULL
		ORDER BY last_updated DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := m.DB.QueryContext(ctx, query, inventoryID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*InventoryProduct
	for rows.Next() {
		var ip InventoryProduct
		if err := rows.Scan(
			&ip.ID, &ip.InventoryID, &ip.ProductVariantID, &ip.Quantity, &ip.Unit, &ip.CreatedAt, &ip.LastUpdated, &ip.DeletedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, &ip)
	}
	return products, rows.Err()
}
