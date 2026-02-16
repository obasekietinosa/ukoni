package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"ukoni/internal/database"
)

type Plan struct {
	ID           string     `json:"id"`
	InventoryID  string     `json:"inventory_id"`
	ParentPlanID *string    `json:"parent_plan_id,omitempty"`
	Title        string     `json:"title"`
	Description  *string    `json:"description,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

type PlanItem struct {
	ID         string     `json:"id"`
	PlanID     string     `json:"plan_id"`
	TargetType string     `json:"target_type"`
	TargetID   string     `json:"target_id"`
	Quantity   *float64   `json:"quantity,omitempty"`
	Unit       *string    `json:"unit,omitempty"`
	Note       *string    `json:"note,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`

	// Join fields
	CanonicalProduct *CanonicalProduct `json:"canonical_product,omitempty"`
	ProductVariant   *ProductVariant   `json:"product_variant,omitempty"`
	Product          *Product          `json:"product,omitempty"`
}

type PlanShoppingList struct {
	PlanID         string    `json:"plan_id"`
	ShoppingListID string    `json:"shopping_list_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type PlanModel struct {
	DB *sql.DB
}

func (m *PlanModel) Create(ctx context.Context, dbtx database.DBTX, plan *Plan) error {
	query := `
		INSERT INTO plans (inventory_id, parent_plan_id, title, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return dbtx.QueryRowContext(ctx, query,
		plan.InventoryID,
		plan.ParentPlanID,
		plan.Title,
		plan.Description,
	).Scan(&plan.ID, &plan.CreatedAt, &plan.UpdatedAt)
}

func (m *PlanModel) GetByID(ctx context.Context, id string) (*Plan, error) {
	query := `
		SELECT id, inventory_id, parent_plan_id, title, description, created_at, updated_at, deleted_at
		FROM plans
		WHERE id = $1 AND deleted_at IS NULL
	`
	var p Plan
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.InventoryID, &p.ParentPlanID, &p.Title, &p.Description, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (m *PlanModel) List(ctx context.Context, inventoryID string, limit, offset int, parentPlanID *string) ([]*Plan, error) {
	query := `
		SELECT id, inventory_id, parent_plan_id, title, description, created_at, updated_at, deleted_at
		FROM plans
		WHERE inventory_id = $1 AND deleted_at IS NULL
	`
	args := []interface{}{inventoryID}
	argCount := 2

	if parentPlanID != nil {
		query += fmt.Sprintf(" AND parent_plan_id = $%d", argCount)
		args = append(args, *parentPlanID)
		argCount++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plans := []*Plan{}
	for rows.Next() {
		var p Plan
		if err := rows.Scan(
			&p.ID, &p.InventoryID, &p.ParentPlanID, &p.Title, &p.Description, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
		); err != nil {
			return nil, err
		}
		plans = append(plans, &p)
	}
	return plans, rows.Err()
}

func (m *PlanModel) Update(ctx context.Context, dbtx database.DBTX, plan *Plan) error {
	query := `
		UPDATE plans
		SET title = $1, description = $2, parent_plan_id = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND deleted_at IS NULL
		RETURNING updated_at
	`
	return dbtx.QueryRowContext(ctx, query,
		plan.Title,
		plan.Description,
		plan.ParentPlanID,
		plan.ID,
	).Scan(&plan.UpdatedAt)
}

func (m *PlanModel) Delete(ctx context.Context, dbtx database.DBTX, id string) error {
	query := `
		UPDATE plans
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

func (m *PlanModel) AddItem(ctx context.Context, dbtx database.DBTX, item *PlanItem) error {
	query := `
		INSERT INTO plan_items (plan_id, target_type, target_id, quantity, unit, note)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return dbtx.QueryRowContext(ctx, query,
		item.PlanID,
		item.TargetType,
		item.TargetID,
		item.Quantity,
		item.Unit,
		item.Note,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
}

func (m *PlanModel) UpdateItem(ctx context.Context, dbtx database.DBTX, item *PlanItem) error {
	query := `
		UPDATE plan_items
		SET quantity = $1, unit = $2, note = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND deleted_at IS NULL
		RETURNING updated_at
	`
	return dbtx.QueryRowContext(ctx, query,
		item.Quantity,
		item.Unit,
		item.Note,
		item.ID,
	).Scan(&item.UpdatedAt)
}

func (m *PlanModel) RemoveItem(ctx context.Context, dbtx database.DBTX, id string) error {
	query := `
		UPDATE plan_items
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

func (m *PlanModel) GetItems(ctx context.Context, planID string) ([]*PlanItem, error) {
	query := `
		SELECT
			pi.id, pi.plan_id, pi.target_type, pi.target_id, pi.quantity, pi.unit, pi.note, pi.created_at, pi.updated_at, pi.deleted_at,
			cp.id, cp.name, cp.category_id,
			pv.id, pv.product_id, pv.variant_name, pv.sku, pv.unit, pv.size,
			p.id, p.name, p.brand,
			p_direct.id, p_direct.name, p_direct.brand
		FROM plan_items pi
		LEFT JOIN canonical_products cp ON pi.target_type = 'canonical_product' AND pi.target_id = cp.id
		LEFT JOIN product_variants pv ON pi.target_type = 'product_variant' AND pi.target_id = pv.id
		LEFT JOIN products p ON pv.product_id = p.id
		LEFT JOIN products p_direct ON pi.target_type = 'product' AND pi.target_id = p_direct.id
		WHERE pi.plan_id = $1 AND pi.deleted_at IS NULL
		ORDER BY pi.created_at ASC
	`
	rows, err := m.DB.QueryContext(ctx, query, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*PlanItem{}
	for rows.Next() {
		var i PlanItem
		var cpID, cpCategory *string
		var cpName *string

		var pvID, pvProdID *string
		var pvName, pvSku, pvUnit *string
		var pvSize *float64

		var pID, pName, pBrand *string
		var pdID, pdName, pdBrand *string

		err := rows.Scan(
			&i.ID, &i.PlanID, &i.TargetType, &i.TargetID, &i.Quantity, &i.Unit, &i.Note, &i.CreatedAt, &i.UpdatedAt, &i.DeletedAt,
			&cpID, &cpName, &cpCategory,
			&pvID, &pvProdID, &pvName, &pvSku, &pvUnit, &pvSize,
			&pID, &pName, &pBrand,
			&pdID, &pdName, &pdBrand,
		)
		if err != nil {
			return nil, err
		}

		if i.TargetType == "canonical_product" && cpID != nil {
			i.CanonicalProduct = &CanonicalProduct{
				ID:   *cpID,
				Name: *cpName,
			}
			if cpCategory != nil {
				i.CanonicalProduct.CategoryID = cpCategory
			}
		} else if i.TargetType == "product_variant" && pvID != nil {
			i.ProductVariant = &ProductVariant{
				ID:          *pvID,
				ProductID:   *pvProdID,
				VariantName: *pvName,
			}
			if pvSku != nil {
				i.ProductVariant.SKU = pvSku
			}
			if pvUnit != nil {
				i.ProductVariant.Unit = pvUnit
			}
			if pvSize != nil {
				i.ProductVariant.Size = pvSize
			}

			if pID != nil {
				i.Product = &Product{
					ID:   *pID,
					Name: *pName,
				}
				if pBrand != nil {
					i.Product.Brand = pBrand
				}
			}
		} else if i.TargetType == "product" && pdID != nil {
			i.Product = &Product{
				ID:   *pdID,
				Name: *pdName,
			}
			if pdBrand != nil {
				i.Product.Brand = pdBrand
			}
		}

		items = append(items, &i)
	}
	return items, rows.Err()
}

func (m *PlanModel) GetItem(ctx context.Context, id string) (*PlanItem, error) {
	query := `
		SELECT id, plan_id, target_type, target_id, quantity, unit, note, created_at, updated_at, deleted_at
		FROM plan_items
		WHERE id = $1 AND deleted_at IS NULL
	`
	var i PlanItem
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&i.ID, &i.PlanID, &i.TargetType, &i.TargetID, &i.Quantity, &i.Unit, &i.Note, &i.CreatedAt, &i.UpdatedAt, &i.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &i, nil
}

func (m *PlanModel) LinkShoppingList(ctx context.Context, dbtx database.DBTX, planID, shoppingListID string) error {
	query := `
		INSERT INTO plan_shopping_lists (plan_id, shopping_list_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`
	_, err := dbtx.ExecContext(ctx, query, planID, shoppingListID)
	return err
}

func (m *PlanModel) UnlinkShoppingList(ctx context.Context, dbtx database.DBTX, planID, shoppingListID string) error {
	query := `
		DELETE FROM plan_shopping_lists
		WHERE plan_id = $1 AND shopping_list_id = $2
	`
	_, err := dbtx.ExecContext(ctx, query, planID, shoppingListID)
	return err
}

func (m *PlanModel) GetShoppingLists(ctx context.Context, planID string) ([]string, error) {
	query := `
		SELECT shopping_list_id
		FROM plan_shopping_lists
		WHERE plan_id = $1
	`
	rows, err := m.DB.QueryContext(ctx, query, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// Helper to get children (sub-plans)
func (m *PlanModel) GetChildren(ctx context.Context, planID string) ([]*Plan, error) {
	query := `
		SELECT id, inventory_id, parent_plan_id, title, description, created_at, updated_at, deleted_at
		FROM plans
		WHERE parent_plan_id = $1 AND deleted_at IS NULL
		ORDER BY created_at ASC
	`
	rows, err := m.DB.QueryContext(ctx, query, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plans := []*Plan{}
	for rows.Next() {
		var p Plan
		if err := rows.Scan(
			&p.ID, &p.InventoryID, &p.ParentPlanID, &p.Title, &p.Description, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
		); err != nil {
			return nil, err
		}
		plans = append(plans, &p)
	}
	return plans, rows.Err()
}
