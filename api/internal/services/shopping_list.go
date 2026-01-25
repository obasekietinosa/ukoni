package services

import (
	"context"
	"errors"
	"ukoni/internal/models"
)

type ShoppingListService struct {
	ShoppingListModel  *models.ShoppingListModel
	InventoryModel     *models.InventoryModel
	MembershipModel    *models.MembershipModel
	ActivityLogService *ActivityLogService
}

func (s *ShoppingListService) checkAccess(ctx context.Context, userID, inventoryID string) (bool, error) {
	// Check if owner
	inv, err := s.InventoryModel.GetByID(inventoryID)
	if err != nil {
		return false, err
	}
	if inv.OwnerUserID == userID {
		return true, nil
	}

	// Check if member
	_, err = s.MembershipModel.GetMembership(inventoryID, userID)
	if err == nil {
		return true, nil
	}
	return false, nil
}

func (s *ShoppingListService) CreateList(ctx context.Context, userID, inventoryID, name string) (*models.ShoppingList, error) {
	// Verify access
	hasAccess, err := s.checkAccess(ctx, userID, inventoryID)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, errors.New("unauthorized")
	}

	list := &models.ShoppingList{
		InventoryID: inventoryID,
		Name:        name,
		CreatedBy:   userID,
	}

	if err := s.ShoppingListModel.CreateList(ctx, list); err != nil {
		return nil, err
	}

	if s.ActivityLogService != nil {
		s.ActivityLogService.LogActivity(ctx, s.ShoppingListModel.DB, &inventoryID, &userID, "shopping_list.created", "shopping_list", &list.ID, nil)
	}

	return list, nil
}

func (s *ShoppingListService) ListLists(ctx context.Context, userID, inventoryID string) ([]*models.ShoppingList, error) {
	hasAccess, err := s.checkAccess(ctx, userID, inventoryID)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, errors.New("unauthorized")
	}

	return s.ShoppingListModel.ListLists(ctx, inventoryID)
}

func (s *ShoppingListService) GetList(ctx context.Context, userID, listID string) (*models.ShoppingList, error) {
	list, err := s.ShoppingListModel.GetList(ctx, listID)
	if err != nil {
		return nil, err
	}

	hasAccess, err := s.checkAccess(ctx, userID, list.InventoryID)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, errors.New("unauthorized")
	}

	return list, nil
}

func (s *ShoppingListService) UpdateList(ctx context.Context, userID, listID, name string) (*models.ShoppingList, error) {
	list, err := s.GetList(ctx, userID, listID)
	if err != nil {
		return nil, err
	}

	list.Name = name
	if err := s.ShoppingListModel.UpdateList(ctx, list); err != nil {
		return nil, err
	}

	if s.ActivityLogService != nil {
		s.ActivityLogService.LogActivity(ctx, s.ShoppingListModel.DB, &list.InventoryID, &userID, "shopping_list.updated", "shopping_list", &list.ID, nil)
	}

	return list, nil
}

func (s *ShoppingListService) DeleteList(ctx context.Context, userID, listID string) error {
	list, err := s.GetList(ctx, userID, listID)
	if err != nil {
		return err
	}

	if err := s.ShoppingListModel.DeleteList(ctx, listID); err != nil {
		return err
	}

	if s.ActivityLogService != nil {
		s.ActivityLogService.LogActivity(ctx, s.ShoppingListModel.DB, &list.InventoryID, &userID, "shopping_list.deleted", "shopping_list", &list.ID, nil)
	}

	return nil
}

func (s *ShoppingListService) AddItem(ctx context.Context, userID, listID string, item *models.ShoppingListItem) (*models.ShoppingListItem, error) {
	list, err := s.GetList(ctx, userID, listID)
	if err != nil {
		return nil, err
	}

	item.ShoppingListID = listID
	if err := s.ShoppingListModel.AddItem(ctx, item); err != nil {
		return nil, err
	}

	if s.ActivityLogService != nil {
		s.ActivityLogService.LogActivity(ctx, s.ShoppingListModel.DB, &list.InventoryID, &userID, "shopping_list_item.created", "shopping_list_item", &item.ID, nil)
	}

	return item, nil
}

func (s *ShoppingListService) ListItems(ctx context.Context, userID, listID string) ([]*models.ShoppingListItem, error) {
	_, err := s.GetList(ctx, userID, listID) // check access
	if err != nil {
		return nil, err
	}

	return s.ShoppingListModel.ListItems(ctx, listID)
}

func (s *ShoppingListService) UpdateItem(ctx context.Context, userID, itemID string, notes *string, preferredOutletID *string) (*models.ShoppingListItem, error) {
	item, err := s.ShoppingListModel.GetItem(ctx, itemID)
	if err != nil {
		return nil, err
	}

	list, err := s.GetList(ctx, userID, item.ShoppingListID)
	if err != nil {
		return nil, err
	}

	if notes != nil {
		item.Notes = notes
	}
	if preferredOutletID != nil {
		item.PreferredOutletID = preferredOutletID
	}

	if err := s.ShoppingListModel.UpdateItem(ctx, item); err != nil {
		return nil, err
	}

	if s.ActivityLogService != nil {
		s.ActivityLogService.LogActivity(ctx, s.ShoppingListModel.DB, &list.InventoryID, &userID, "shopping_list_item.updated", "shopping_list_item", &item.ID, nil)
	}

	return item, nil
}

func (s *ShoppingListService) DeleteItem(ctx context.Context, userID, itemID string) error {
	item, err := s.ShoppingListModel.GetItem(ctx, itemID)
	if err != nil {
		return err
	}

	list, err := s.GetList(ctx, userID, item.ShoppingListID)
	if err != nil {
		return err
	}

	if err := s.ShoppingListModel.DeleteItem(ctx, itemID); err != nil {
		return err
	}

	if s.ActivityLogService != nil {
		s.ActivityLogService.LogActivity(ctx, s.ShoppingListModel.DB, &list.InventoryID, &userID, "shopping_list_item.deleted", "shopping_list_item", &item.ID, nil)
	}

	return nil
}
