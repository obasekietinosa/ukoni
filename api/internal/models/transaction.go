package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
	"ukoni/internal/database"
)

type Transaction struct {
	ID              string     `json:"id"`
	InventoryID     string     `json:"inventory_id"`
	OutletID        *string    `json:"outlet_id,omitempty"`
	CreatedByUserID string     `json:"created_by_user_id"`
	TransactionDate time.Time  `json:"transaction_date"`
	TotalAmount     *float64   `json:"total_amount,omitempty"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

type TransactionItem struct {
	ID                 string     `json:"id"`
	TransactionID      string     `json:"transaction_id"`
	ProductVariantID   string     `json:"product_variant_id"`
	Quantity           float64    `json:"quantity"`
	PricePerUnit       *float64   `json:"price_per_unit,omitempty"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
	ShoppingListItemID *string    `json:"shopping_list_item_id,omitempty"`
}

type TransactionModel struct {
	DB *sql.DB
}

func (m *TransactionModel) Create(ctx context.Context, dbtx database.DBTX, t *Transaction) error {
	query := `
		INSERT INTO transactions (inventory_id, outlet_id, created_by_user_id, transaction_date, total_amount)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, transaction_date
	`
	return dbtx.QueryRowContext(ctx, query,
		t.InventoryID,
		t.OutletID,
		t.CreatedByUserID,
		t.TransactionDate,
		t.TotalAmount,
	).Scan(&t.ID, &t.TransactionDate)
}

func (m *TransactionModel) CreateItem(ctx context.Context, dbtx database.DBTX, item *TransactionItem) error {
	query := `
		INSERT INTO transaction_items (transaction_id, product_variant_id, quantity, price_per_unit, shopping_list_item_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return dbtx.QueryRowContext(ctx, query,
		item.TransactionID,
		item.ProductVariantID,
		item.Quantity,
		item.PricePerUnit,
		item.ShoppingListItemID,
	).Scan(&item.ID)
}

func (m *TransactionModel) CreateItems(ctx context.Context, dbtx database.DBTX, items []*TransactionItem) error {
	if len(items) == 0 {
		return nil
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString("INSERT INTO transaction_items (transaction_id, product_variant_id, quantity, price_per_unit, shopping_list_item_id) VALUES ")

	args := make([]interface{}, 0, len(items)*5)

	for i, item := range items {
		n := i * 5
		if i > 0 {
			queryBuilder.WriteString(",")
		}
		fmt.Fprintf(&queryBuilder, "($%d, $%d, $%d, $%d, $%d)", n+1, n+2, n+3, n+4, n+5)
		args = append(args, item.TransactionID, item.ProductVariantID, item.Quantity, item.PricePerUnit, item.ShoppingListItemID)
	}

	queryBuilder.WriteString(" RETURNING id")

	rows, err := dbtx.QueryContext(ctx, queryBuilder.String(), args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		if i >= len(items) {
			break
		}
		if err := rows.Scan(&items[i].ID); err != nil {
			return err
		}
		i++
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func (m *TransactionModel) GetByID(ctx context.Context, id string) (*Transaction, error) {
	query := `
		SELECT id, inventory_id, outlet_id, created_by_user_id, transaction_date, total_amount, deleted_at
		FROM transactions
		WHERE id = $1 AND deleted_at IS NULL
	`
	var t Transaction
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&t.ID,
		&t.InventoryID,
		&t.OutletID,
		&t.CreatedByUserID,
		&t.TransactionDate,
		&t.TotalAmount,
		&t.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (m *TransactionModel) ListByInventory(ctx context.Context, inventoryID string, limit, offset int) ([]*Transaction, error) {
	query := `
		SELECT id, inventory_id, outlet_id, created_by_user_id, transaction_date, total_amount, deleted_at
		FROM transactions
		WHERE inventory_id = $1 AND deleted_at IS NULL
		ORDER BY transaction_date DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := m.DB.QueryContext(ctx, query, inventoryID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(
			&t.ID,
			&t.InventoryID,
			&t.OutletID,
			&t.CreatedByUserID,
			&t.TransactionDate,
			&t.TotalAmount,
			&t.DeletedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return transactions, nil
}

func (m *TransactionModel) GetItems(ctx context.Context, transactionID string) ([]*TransactionItem, error) {
	query := `
		SELECT id, transaction_id, product_variant_id, quantity, price_per_unit, shopping_list_item_id, deleted_at
		FROM transaction_items
		WHERE transaction_id = $1 AND deleted_at IS NULL
	`
	rows, err := m.DB.QueryContext(ctx, query, transactionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*TransactionItem
	for rows.Next() {
		var item TransactionItem
		if err := rows.Scan(
			&item.ID,
			&item.TransactionID,
			&item.ProductVariantID,
			&item.Quantity,
			&item.PricePerUnit,
			&item.ShoppingListItemID,
			&item.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
