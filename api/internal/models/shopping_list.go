package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ShoppingList struct {
	ID            string     `json:"id"`
	InventoryID   string     `json:"inventory_id"`
	Name          string     `json:"name"`
	CreatedBy     string     `json:"created_by"`
	CreatedAt     time.Time  `json:"created_at"`
	LastUpdatedAt time.Time  `json:"last_updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

type ShoppingListItem struct {
	ID                string     `json:"id"`
	ShoppingListID    string     `json:"shopping_list_id"`
	TargetType        string     `json:"target_type"` // 'canonical_product' or 'product_variant'
	TargetID          string     `json:"target_id"`
	PreferredOutletID *string    `json:"preferred_outlet_id,omitempty"`
	Notes             *string    `json:"notes,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`

	// Join fields
	CanonicalProduct *CanonicalProduct `json:"canonical_product,omitempty"`
	ProductVariant   *ProductVariant   `json:"product_variant,omitempty"`
	PreferredOutlet  *Outlet           `json:"preferred_outlet,omitempty"`
}

type ShoppingListRepository interface {
	CreateList(ctx context.Context, list *ShoppingList) error
	GetList(ctx context.Context, id string) (*ShoppingList, error)
	ListLists(ctx context.Context, inventoryID string) ([]*ShoppingList, error)
	UpdateList(ctx context.Context, list *ShoppingList) error
	DeleteList(ctx context.Context, id string) error

	AddItem(ctx context.Context, item *ShoppingListItem) error
	GetItem(ctx context.Context, id string) (*ShoppingListItem, error)
	ListItems(ctx context.Context, listID string) ([]*ShoppingListItem, error)
	UpdateItem(ctx context.Context, item *ShoppingListItem) error
	DeleteItem(ctx context.Context, id string) error
}

type ShoppingListModel struct {
	DB *sql.DB
}

func (m *ShoppingListModel) CreateList(ctx context.Context, list *ShoppingList) error {
	query := `
		INSERT INTO shopping_lists (name, inventory_id, created_by)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, last_updated_at, deleted_at
	`
	return m.DB.QueryRowContext(ctx, query, list.Name, list.InventoryID, list.CreatedBy).Scan(
		&list.ID, &list.CreatedAt, &list.LastUpdatedAt, &list.DeletedAt,
	)
}

func (m *ShoppingListModel) GetList(ctx context.Context, id string) (*ShoppingList, error) {
	query := `
		SELECT id, inventory_id, name, created_by, created_at, last_updated_at, deleted_at
		FROM shopping_lists
		WHERE id = $1 AND deleted_at IS NULL
	`
	var list ShoppingList
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&list.ID, &list.InventoryID, &list.Name, &list.CreatedBy,
		&list.CreatedAt, &list.LastUpdatedAt, &list.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func (m *ShoppingListModel) ListLists(ctx context.Context, inventoryID string) ([]*ShoppingList, error) {
	query := `
		SELECT id, inventory_id, name, created_by, created_at, last_updated_at, deleted_at
		FROM shopping_lists
		WHERE inventory_id = $1 AND deleted_at IS NULL
		ORDER BY last_updated_at DESC
	`
	rows, err := m.DB.QueryContext(ctx, query, inventoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lists []*ShoppingList
	for rows.Next() {
		var list ShoppingList
		if err := rows.Scan(
			&list.ID, &list.InventoryID, &list.Name, &list.CreatedBy,
			&list.CreatedAt, &list.LastUpdatedAt, &list.DeletedAt,
		); err != nil {
			return nil, err
		}
		lists = append(lists, &list)
	}
	return lists, nil
}

func (m *ShoppingListModel) UpdateList(ctx context.Context, list *ShoppingList) error {
	query := `
		UPDATE shopping_lists
		SET name = $1, last_updated_at = COALESCE($2, CURRENT_TIMESTAMP)
		WHERE id = $3 AND deleted_at IS NULL
		RETURNING last_updated_at
	`
	return m.DB.QueryRowContext(ctx, query, list.Name, list.LastUpdatedAt, list.ID).Scan(&list.LastUpdatedAt)
}

func (m *ShoppingListModel) DeleteList(ctx context.Context, id string) error {
	query := `
		UPDATE shopping_lists
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func (m *ShoppingListModel) AddItem(ctx context.Context, item *ShoppingListItem) error {
	query := `
		INSERT INTO shopping_list_items (shopping_list_id, target_type, target_id, preferred_outlet_id, notes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return m.DB.QueryRowContext(ctx, query,
		item.ShoppingListID, item.TargetType, item.TargetID, item.PreferredOutletID, item.Notes,
	).Scan(&item.ID, &item.CreatedAt)
}

func (m *ShoppingListModel) GetItem(ctx context.Context, id string) (*ShoppingListItem, error) {
	query := `
		SELECT id, shopping_list_id, target_type, target_id, preferred_outlet_id, notes, created_at, deleted_at
		FROM shopping_list_items
		WHERE id = $1 AND deleted_at IS NULL
	`
	var item ShoppingListItem
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&item.ID, &item.ShoppingListID, &item.TargetType, &item.TargetID,
		&item.PreferredOutletID, &item.Notes, &item.CreatedAt, &item.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (m *ShoppingListModel) ListItems(ctx context.Context, listID string) ([]*ShoppingListItem, error) {
	query := `
		SELECT 
			sli.id, sli.shopping_list_id, sli.target_type, sli.target_id, sli.preferred_outlet_id, sli.notes, sli.created_at, sli.deleted_at,
			cp.id, cp.name, cp.category_id,
			pv.id, pv.product_id, pv.variant_name, pv.sku, pv.unit, pv.size,
			p.id, p.name, p.brand,
			o.id, o.name, o.address
		FROM shopping_list_items sli
		LEFT JOIN canonical_products cp ON sli.target_type = 'canonical_product' AND sli.target_id = cp.id
		LEFT JOIN product_variants pv ON sli.target_type = 'product_variant' AND sli.target_id = pv.id
		LEFT JOIN products p ON pv.product_id = p.id
		LEFT JOIN outlets o ON sli.preferred_outlet_id = o.id
		WHERE sli.shopping_list_id = $1 AND sli.deleted_at IS NULL
		ORDER BY sli.created_at ASC
	`

	rows, err := m.DB.QueryContext(ctx, query, listID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*ShoppingListItem
	for rows.Next() {
		var item ShoppingListItem
		var cpID, cpCategory *string
		var cpName *string

		var pvID, pvProdID *string
		var pvName, pvSku, pvUnit *string
		var pvSize *float64

		var pID *string
		var pName, pBrand *string

		var oID *string
		var oName, oAddress *string

		err := rows.Scan(
			&item.ID, &item.ShoppingListID, &item.TargetType, &item.TargetID, &item.PreferredOutletID, &item.Notes, &item.CreatedAt, &item.DeletedAt,
			&cpID, &cpName, &cpCategory,
			&pvID, &pvProdID, &pvName, &pvSku, &pvUnit, &pvSize,
			&pID, &pName, &pBrand,
			&oID, &oName, &oAddress,
		)
		if err != nil {
			return nil, err
		}

		if item.TargetType == "canonical_product" && cpID != nil {
			item.CanonicalProduct = &CanonicalProduct{
				ID:   *cpID,
				Name: *cpName,
			}
			if cpCategory != nil {
				item.CanonicalProduct.CategoryID = cpCategory // Assuming CategoryID is string pointer
			}
		} else if item.TargetType == "product_variant" && pvID != nil {
			item.ProductVariant = &ProductVariant{
				ID:          *pvID,
				ProductID:   *pvProdID,
				VariantName: *pvName,
			}
			if pvSku != nil {
				item.ProductVariant.SKU = pvSku
			}
			if pvUnit != nil {
				item.ProductVariant.Unit = pvUnit
			}
			if pvSize != nil {
				item.ProductVariant.Size = pvSize
			}
		}

		if oID != nil {
			outletID, err := uuid.Parse(*oID)
			if err == nil {
				item.PreferredOutlet = &Outlet{
					ID:   outletID,
					Name: *oName,
				}
				if oAddress != nil {
					item.PreferredOutlet.Address = *oAddress
				}
			}
		}

		items = append(items, &item)
	}
	return items, nil
}

func (m *ShoppingListModel) UpdateItem(ctx context.Context, item *ShoppingListItem) error {
	query := `
		UPDATE shopping_list_items
		SET notes = $1, preferred_outlet_id = $2
		WHERE id = $3 AND deleted_at IS NULL
	`
	_, err := m.DB.ExecContext(ctx, query, item.Notes, item.PreferredOutletID, item.ID)
	return err
}

func (m *ShoppingListModel) DeleteItem(ctx context.Context, id string) error {
	query := `
		UPDATE shopping_list_items
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
