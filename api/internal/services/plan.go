package services

import (
	"context"
	"database/sql"
	"fmt"
	"ukoni/internal/models"
)

type PlanService struct {
	DB                    *sql.DB
	PlanModel             *models.PlanModel
	InventoryModel        *models.InventoryModel
	ProductModel          *models.ProductModel
	CanonicalProductModel *models.CanonicalProductModel
	ShoppingListModel     *models.ShoppingListModel
}

func (s *PlanService) CreatePlan(ctx context.Context, inventoryID string, parentPlanID *string, title, description string) (*models.Plan, error) {
	if inventoryID == "" {
		return nil, fmt.Errorf("%w: inventory id is required", ErrInvalidInput)
	}
	if title == "" {
		return nil, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	// Validate inventory exists
	inventory, err := s.InventoryModel.GetByID(inventoryID)
	if err != nil {
		return nil, err
	}
	if inventory == nil {
		return nil, fmt.Errorf("%w: inventory not found", ErrInvalidInput)
	}

	// Validate parent plan if provided
	if parentPlanID != nil {
		parent, err := s.PlanModel.GetByID(ctx, *parentPlanID)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, fmt.Errorf("%w: parent plan not found", ErrInvalidInput)
		}
		if parent.InventoryID != inventoryID {
			return nil, fmt.Errorf("%w: parent plan belongs to different inventory", ErrInvalidInput)
		}
	}

	plan := &models.Plan{
		InventoryID:  inventoryID,
		ParentPlanID: parentPlanID,
		Title:        title,
	}
	if description != "" {
		plan.Description = &description
	}

	err = s.PlanModel.Create(ctx, s.DB, plan)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

type PlanWithDetails struct {
	*models.Plan
	Items         []*models.PlanItem `json:"items"`
	Children      []*models.Plan     `json:"children"`
	ShoppingLists []string           `json:"shopping_lists"` // Just IDs for now
}

func (s *PlanService) GetPlan(ctx context.Context, id string) (*PlanWithDetails, error) {
	plan, err := s.PlanModel.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, nil
	}

	items, err := s.PlanModel.GetItems(ctx, id)
	if err != nil {
		return nil, err
	}

	children, err := s.PlanModel.GetChildren(ctx, id)
	if err != nil {
		return nil, err
	}

	shoppingLists, err := s.PlanModel.GetShoppingLists(ctx, id)
	if err != nil {
		return nil, err
	}

	return &PlanWithDetails{
		Plan:          plan,
		Items:         items,
		Children:      children,
		ShoppingLists: shoppingLists,
	}, nil
}

func (s *PlanService) ListPlans(ctx context.Context, inventoryID string, limit, offset int, parentPlanID *string) ([]*models.Plan, error) {
	if inventoryID == "" {
		return nil, fmt.Errorf("%w: inventory id is required", ErrInvalidInput)
	}
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.PlanModel.List(ctx, inventoryID, limit, offset, parentPlanID)
}

func (s *PlanService) UpdatePlan(ctx context.Context, id, title, description string, parentPlanID *string) (*models.Plan, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: plan id is required", ErrInvalidInput)
	}
	if title == "" {
		return nil, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	existing, err := s.PlanModel.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrNotFound
	}

	// Validate parent plan if provided and changed
	if parentPlanID != nil && (existing.ParentPlanID == nil || *parentPlanID != *existing.ParentPlanID) {
		// Prevent circular dependency (basic check: parent cannot be self)
		if *parentPlanID == id {
			return nil, fmt.Errorf("%w: plan cannot be its own parent", ErrInvalidInput)
		}
		// Fetch parent
		parent, err := s.PlanModel.GetByID(ctx, *parentPlanID)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, fmt.Errorf("%w: parent plan not found", ErrInvalidInput)
		}
		if parent.InventoryID != existing.InventoryID {
			return nil, fmt.Errorf("%w: parent plan belongs to different inventory", ErrInvalidInput)
		}
	}

	plan := &models.Plan{
		ID:           id,
		Title:        title,
		ParentPlanID: parentPlanID,
	}
	if description != "" {
		plan.Description = &description
	}

	err = s.PlanModel.Update(ctx, s.DB, plan)
	if err != nil {
		return nil, err
	}

	// Reload to get full state
	return s.PlanModel.GetByID(ctx, id)
}

func (s *PlanService) DeletePlan(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("%w: plan id is required", ErrInvalidInput)
	}
	err := s.PlanModel.Delete(ctx, s.DB, id)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}
	return err
}

func (s *PlanService) AddItem(ctx context.Context, planID, targetType, targetID string, quantity *float64, unit, note string) (*models.PlanItem, error) {
	if planID == "" {
		return nil, fmt.Errorf("%w: plan id is required", ErrInvalidInput)
	}
	if targetID == "" {
		return nil, fmt.Errorf("%w: target id is required", ErrInvalidInput)
	}
	if targetType != "canonical_product" && targetType != "product" && targetType != "product_variant" {
		return nil, fmt.Errorf("%w: invalid target type", ErrInvalidInput)
	}

	// Validate plan exists
	plan, err := s.PlanModel.GetByID(ctx, planID)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, ErrNotFound
	}

	// Validate target exists
	switch targetType {
	case "canonical_product":
		cp, err := s.CanonicalProductModel.GetByID(ctx, targetID)
		if err != nil {
			return nil, err
		}
		if cp == nil {
			return nil, fmt.Errorf("%w: canonical product not found", ErrInvalidInput)
		}
		if cp.InventoryID != plan.InventoryID {
			return nil, fmt.Errorf("%w: canonical product belongs to different inventory", ErrInvalidInput)
		}
	case "product":
		p, err := s.ProductModel.GetByID(ctx, targetID)
		if err != nil {
			return nil, err
		}
		if p == nil {
			return nil, fmt.Errorf("%w: product not found", ErrInvalidInput)
		}
		if p.InventoryID != plan.InventoryID {
			return nil, fmt.Errorf("%w: product belongs to different inventory", ErrInvalidInput)
		}
	case "product_variant":
		v, err := s.ProductModel.GetVariant(ctx, targetID)
		if err != nil {
			return nil, err
		}
		if v == nil {
			return nil, fmt.Errorf("%w: product variant not found", ErrInvalidInput)
		}
		// Need to check if variant's product belongs to same inventory
		p, err := s.ProductModel.GetByID(ctx, v.ProductID)
		if err != nil {
			return nil, err
		}
		if p == nil {
			// Inconsistent state
			return nil, fmt.Errorf("%w: product for variant not found", ErrInvalidInput)
		}
		if p.InventoryID != plan.InventoryID {
			return nil, fmt.Errorf("%w: product variant belongs to different inventory", ErrInvalidInput)
		}
	}

	item := &models.PlanItem{
		PlanID:     planID,
		TargetType: targetType,
		TargetID:   targetID,
		Quantity:   quantity,
	}
	if unit != "" {
		item.Unit = &unit
	}
	if note != "" {
		item.Note = &note
	}

	err = s.PlanModel.AddItem(ctx, s.DB, item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *PlanService) UpdateItem(ctx context.Context, id string, quantity *float64, unit, note string) (*models.PlanItem, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: item id is required", ErrInvalidInput)
	}

	item := &models.PlanItem{
		ID:       id,
		Quantity: quantity,
	}
	if unit != "" {
		item.Unit = &unit
	}
	if note != "" {
		item.Note = &note
	}

	err := s.PlanModel.UpdateItem(ctx, s.DB, item)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return s.PlanModel.GetItem(ctx, id)
}

func (s *PlanService) GetItem(ctx context.Context, id string) (*models.PlanItem, error) {
	return s.PlanModel.GetItem(ctx, id)
}

func (s *PlanService) RemoveItem(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("%w: item id is required", ErrInvalidInput)
	}
	err := s.PlanModel.RemoveItem(ctx, s.DB, id)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}
	return err
}

func (s *PlanService) LinkShoppingList(ctx context.Context, planID, shoppingListID string) error {
	if planID == "" || shoppingListID == "" {
		return fmt.Errorf("%w: plan id and shopping list id are required", ErrInvalidInput)
	}

	plan, err := s.PlanModel.GetByID(ctx, planID)
	if err != nil {
		return err
	}
	if plan == nil {
		return ErrNotFound
	}

	sl, err := s.ShoppingListModel.GetList(ctx, shoppingListID)
	if err != nil {
		return err
	}
	if sl == nil {
		return fmt.Errorf("%w: shopping list not found", ErrInvalidInput)
	}

	if plan.InventoryID != sl.InventoryID {
		return fmt.Errorf("%w: inventory mismatch", ErrInvalidInput)
	}

	return s.PlanModel.LinkShoppingList(ctx, s.DB, planID, shoppingListID)
}

func (s *PlanService) UnlinkShoppingList(ctx context.Context, planID, shoppingListID string) error {
	if planID == "" || shoppingListID == "" {
		return fmt.Errorf("%w: plan id and shopping list id are required", ErrInvalidInput)
	}
	return s.PlanModel.UnlinkShoppingList(ctx, s.DB, planID, shoppingListID)
}

func (s *PlanService) CreateShoppingListFromGroup(ctx context.Context, planID, name, userID string) (*models.ShoppingList, error) {
	if planID == "" {
		return nil, fmt.Errorf("%w: plan id is required", ErrInvalidInput)
	}

	plan, err := s.PlanModel.GetByID(ctx, planID)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, ErrNotFound
	}

	if name == "" {
		name = fmt.Sprintf("%s Shopping List", plan.Title)
	}

	// Create shopping list
	sl := &models.ShoppingList{
		InventoryID: plan.InventoryID,
		Name:        name,
		CreatedBy:   userID,
	}
	err = s.ShoppingListModel.CreateList(ctx, sl)
	if err != nil {
		return nil, err
	}

	// Helper to add items from a plan to the shopping list
	addItemsFromPlan := func(pID string) error {
		items, err := s.PlanModel.GetItems(ctx, pID)
		if err != nil {
			return err
		}
		for _, item := range items {
			slItem := &models.ShoppingListItem{
				ShoppingListID: sl.ID,
				TargetType:     item.TargetType,
				TargetID:       item.TargetID,
				Quantity:       item.Quantity,
				Unit:           item.Unit,
				Notes:          item.Note,
			}
			if err := s.ShoppingListModel.AddItem(ctx, slItem); err != nil {
				return err
			}
		}
		return nil
	}

	// Add items from parent plan
	if err := addItemsFromPlan(planID); err != nil {
		return nil, err
	}

	// Add items from children
	children, err := s.PlanModel.GetChildren(ctx, planID)
	if err != nil {
		return nil, err
	}
	for _, child := range children {
		if err := addItemsFromPlan(child.ID); err != nil {
			return nil, err
		}
	}

	// Link shopping list to parent plan
	if err := s.PlanModel.LinkShoppingList(ctx, s.DB, planID, sl.ID); err != nil {
		return nil, err
	}

	return sl, nil
}
