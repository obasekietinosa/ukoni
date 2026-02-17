package models

import (
	"context"
	"database/sql"
	"time"
	"ukoni/internal/database"
)

type PlanGroup struct {
	ID          string     `json:"id"`
	InventoryID string     `json:"inventory_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type PlanGroupModel struct {
	DB *sql.DB
}

func (m *PlanGroupModel) Create(ctx context.Context, dbtx database.DBTX, group *PlanGroup) error {
	query := `
		INSERT INTO plan_groups (inventory_id, title, description)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	return dbtx.QueryRowContext(ctx, query,
		group.InventoryID,
		group.Title,
		group.Description,
	).Scan(&group.ID, &group.CreatedAt, &group.UpdatedAt)
}

func (m *PlanGroupModel) GetByID(ctx context.Context, id string) (*PlanGroup, error) {
	query := `
		SELECT id, inventory_id, title, description, created_at, updated_at, deleted_at
		FROM plan_groups
		WHERE id = $1 AND deleted_at IS NULL
	`
	var g PlanGroup
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&g.ID, &g.InventoryID, &g.Title, &g.Description, &g.CreatedAt, &g.UpdatedAt, &g.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &g, nil
}

func (m *PlanGroupModel) List(ctx context.Context, inventoryID string, limit, offset int) ([]*PlanGroup, error) {
	query := `
		SELECT id, inventory_id, title, description, created_at, updated_at, deleted_at
		FROM plan_groups
		WHERE inventory_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := m.DB.QueryContext(ctx, query, inventoryID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := []*PlanGroup{}
	for rows.Next() {
		var g PlanGroup
		if err := rows.Scan(
			&g.ID, &g.InventoryID, &g.Title, &g.Description, &g.CreatedAt, &g.UpdatedAt, &g.DeletedAt,
		); err != nil {
			return nil, err
		}
		groups = append(groups, &g)
	}
	return groups, rows.Err()
}

func (m *PlanGroupModel) Update(ctx context.Context, dbtx database.DBTX, group *PlanGroup) error {
	query := `
		UPDATE plan_groups
		SET title = $1, description = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND deleted_at IS NULL
		RETURNING updated_at
	`
	return dbtx.QueryRowContext(ctx, query,
		group.Title,
		group.Description,
		group.ID,
	).Scan(&group.UpdatedAt)
}

func (m *PlanGroupModel) Delete(ctx context.Context, dbtx database.DBTX, id string) error {
	query := `
		UPDATE plan_groups
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

func (m *PlanGroupModel) AddPlan(ctx context.Context, dbtx database.DBTX, groupID, planID string) error {
	query := `
		INSERT INTO plan_group_plans (plan_group_id, plan_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`
	_, err := dbtx.ExecContext(ctx, query, groupID, planID)
	return err
}

func (m *PlanGroupModel) RemovePlan(ctx context.Context, dbtx database.DBTX, groupID, planID string) error {
	query := `
		DELETE FROM plan_group_plans
		WHERE plan_group_id = $1 AND plan_id = $2
	`
	_, err := dbtx.ExecContext(ctx, query, groupID, planID)
	return err
}

func (m *PlanGroupModel) GetPlans(ctx context.Context, groupID string) ([]*Plan, error) {
	query := `
		SELECT p.id, p.inventory_id, p.title, p.description, p.created_at, p.updated_at, p.deleted_at
		FROM plans p
		JOIN plan_group_plans pgp ON p.id = pgp.plan_id
		WHERE pgp.plan_group_id = $1 AND p.deleted_at IS NULL
		ORDER BY p.created_at ASC
	`
	rows, err := m.DB.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plans := []*Plan{}
	for rows.Next() {
		var p Plan
		if err := rows.Scan(
			&p.ID, &p.InventoryID, &p.Title, &p.Description, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
		); err != nil {
			return nil, err
		}
		plans = append(plans, &p)
	}
	return plans, rows.Err()
}

func (m *PlanGroupModel) LinkShoppingList(ctx context.Context, dbtx database.DBTX, groupID, shoppingListID string) error {
	query := `
		INSERT INTO plan_group_shopping_lists (plan_group_id, shopping_list_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`
	_, err := dbtx.ExecContext(ctx, query, groupID, shoppingListID)
	return err
}

func (m *PlanGroupModel) UnlinkShoppingList(ctx context.Context, dbtx database.DBTX, groupID, shoppingListID string) error {
	query := `
		DELETE FROM plan_group_shopping_lists
		WHERE plan_group_id = $1 AND shopping_list_id = $2
	`
	_, err := dbtx.ExecContext(ctx, query, groupID, shoppingListID)
	return err
}

func (m *PlanGroupModel) GetShoppingLists(ctx context.Context, groupID string) ([]string, error) {
	query := `
		SELECT shopping_list_id
		FROM plan_group_shopping_lists
		WHERE plan_group_id = $1
	`
	rows, err := m.DB.QueryContext(ctx, query, groupID)
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
