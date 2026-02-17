package services

import (
	"context"
	"database/sql"
	"fmt"
	"ukoni/internal/models"
)

type PlanGroupService struct {
	DB             *sql.DB
	PlanGroupModel *models.PlanGroupModel
	PlanModel      *models.PlanModel
	InventoryModel *models.InventoryModel
	ShoppingListModel *models.ShoppingListModel
}

func (s *PlanGroupService) CreateGroup(ctx context.Context, inventoryID, title, description string) (*models.PlanGroup, error) {
	if inventoryID == "" {
		return nil, fmt.Errorf("%w: inventory id is required", ErrInvalidInput)
	}
	if title == "" {
		return nil, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	inventory, err := s.InventoryModel.GetByID(inventoryID)
	if err != nil {
		return nil, err
	}
	if inventory == nil {
		return nil, fmt.Errorf("%w: inventory not found", ErrInvalidInput)
	}

	group := &models.PlanGroup{
		InventoryID: inventoryID,
		Title:       title,
	}
	if description != "" {
		group.Description = &description
	}

	err = s.PlanGroupModel.Create(ctx, s.DB, group)
	if err != nil {
		return nil, err
	}
	return group, nil
}

type PlanGroupWithDetails struct {
	*models.PlanGroup
	Plans         []*models.Plan `json:"plans"`
	ShoppingLists []string       `json:"shopping_lists"`
}

func (s *PlanGroupService) GetGroup(ctx context.Context, id string) (*PlanGroupWithDetails, error) {
	group, err := s.PlanGroupModel.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, nil
	}

	plans, err := s.PlanGroupModel.GetPlans(ctx, id)
	if err != nil {
		return nil, err
	}

	shoppingLists, err := s.PlanGroupModel.GetShoppingLists(ctx, id)
	if err != nil {
		return nil, err
	}

	return &PlanGroupWithDetails{
		PlanGroup:     group,
		Plans:         plans,
		ShoppingLists: shoppingLists,
	}, nil
}

func (s *PlanGroupService) ListGroups(ctx context.Context, inventoryID string, limit, offset int) ([]*models.PlanGroup, error) {
	if inventoryID == "" {
		return nil, fmt.Errorf("%w: inventory id is required", ErrInvalidInput)
	}
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.PlanGroupModel.List(ctx, inventoryID, limit, offset)
}

func (s *PlanGroupService) UpdateGroup(ctx context.Context, id, title, description string) (*models.PlanGroup, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: group id is required", ErrInvalidInput)
	}
	if title == "" {
		return nil, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	existing, err := s.PlanGroupModel.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrNotFound
	}

	group := &models.PlanGroup{
		ID:    id,
		Title: title,
	}
	if description != "" {
		group.Description = &description
	}

	err = s.PlanGroupModel.Update(ctx, s.DB, group)
	if err != nil {
		return nil, err
	}

	return s.PlanGroupModel.GetByID(ctx, id)
}

func (s *PlanGroupService) DeleteGroup(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("%w: group id is required", ErrInvalidInput)
	}
	err := s.PlanGroupModel.Delete(ctx, s.DB, id)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}
	return err
}

func (s *PlanGroupService) AddPlanToGroup(ctx context.Context, groupID, planID string) error {
	if groupID == "" || planID == "" {
		return fmt.Errorf("%w: group id and plan id are required", ErrInvalidInput)
	}

	group, err := s.PlanGroupModel.GetByID(ctx, groupID)
	if err != nil { return err }
	if group == nil { return ErrNotFound }

	plan, err := s.PlanModel.GetByID(ctx, planID)
	if err != nil { return err }
	if plan == nil { return fmt.Errorf("%w: plan not found", ErrInvalidInput) }

	if group.InventoryID != plan.InventoryID {
		return fmt.Errorf("%w: inventory mismatch", ErrInvalidInput)
	}

	return s.PlanGroupModel.AddPlan(ctx, s.DB, groupID, planID)
}

func (s *PlanGroupService) RemovePlanFromGroup(ctx context.Context, groupID, planID string) error {
	if groupID == "" || planID == "" {
		return fmt.Errorf("%w: group id and plan id are required", ErrInvalidInput)
	}
	return s.PlanGroupModel.RemovePlan(ctx, s.DB, groupID, planID)
}

func (s *PlanGroupService) CreateShoppingListFromGroup(ctx context.Context, groupID, name, userID string) (*models.ShoppingList, error) {
	if groupID == "" {
		return nil, fmt.Errorf("%w: group id is required", ErrInvalidInput)
	}

	group, err := s.PlanGroupModel.GetByID(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, ErrNotFound
	}

	if name == "" {
		name = fmt.Sprintf("%s Shopping List", group.Title)
	}

	sl := &models.ShoppingList{
		InventoryID: group.InventoryID,
		Name:        name,
		CreatedBy:   userID,
	}
	err = s.ShoppingListModel.CreateList(ctx, sl)
	if err != nil {
		return nil, err
	}

	plans, err := s.PlanGroupModel.GetPlans(ctx, groupID)
	if err != nil {
		return nil, err
	}

	for _, p := range plans {
		items, err := s.PlanModel.GetItems(ctx, p.ID)
		if err != nil {
			return nil, err
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
				return nil, err
			}
		}
	}

	if err := s.PlanGroupModel.LinkShoppingList(ctx, s.DB, groupID, sl.ID); err != nil {
		return nil, err
	}

	return sl, nil
}

func (s *PlanGroupService) LinkShoppingList(ctx context.Context, groupID, shoppingListID string) error {
	if groupID == "" || shoppingListID == "" {
		return fmt.Errorf("%w: group id and shopping list id are required", ErrInvalidInput)
	}

	group, err := s.PlanGroupModel.GetByID(ctx, groupID)
	if err != nil {
		return err
	}
	if group == nil {
		return ErrNotFound
	}

	sl, err := s.ShoppingListModel.GetList(ctx, shoppingListID)
	if err != nil {
		return err
	}
	if sl == nil {
		return fmt.Errorf("%w: shopping list not found", ErrInvalidInput)
	}

	if group.InventoryID != sl.InventoryID {
		return fmt.Errorf("%w: inventory mismatch", ErrInvalidInput)
	}

	return s.PlanGroupModel.LinkShoppingList(ctx, s.DB, groupID, shoppingListID)
}

func (s *PlanGroupService) UnlinkShoppingList(ctx context.Context, groupID, shoppingListID string) error {
	if groupID == "" || shoppingListID == "" {
		return fmt.Errorf("%w: group id and shopping list id are required", ErrInvalidInput)
	}
	return s.PlanGroupModel.UnlinkShoppingList(ctx, s.DB, groupID, shoppingListID)
}
