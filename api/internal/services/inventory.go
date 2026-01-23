package services

import (
	"context"
	"database/sql"
	"errors"
	"ukoni/internal/models"
)

type InventoryService struct {
	DB                 *sql.DB
	InventoryModel     *models.InventoryModel
	MembershipModel    *models.MembershipModel
	ActivityLogService *ActivityLogService
}

func (s *InventoryService) CreateInventory(ctx context.Context, userID, name string) (*models.Inventory, error) {
	if name == "" {
		return nil, errors.New("inventory name cannot be empty")
	}

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	inventory := &models.Inventory{
		Name:        name,
		OwnerUserID: userID,
	}

	if err := s.InventoryModel.Create(ctx, tx, inventory); err != nil {
		return nil, err
	}

	// Add owner as an admin member
	err = s.MembershipModel.AddMember(ctx, tx, inventory.ID, userID, "admin")
	if err != nil {
		return nil, err
	}

	if s.ActivityLogService != nil {
		if err := s.ActivityLogService.LogActivity(ctx, tx, &inventory.ID, &userID, "inventory.created", "inventory", &inventory.ID, nil); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return inventory, nil
}

func (s *InventoryService) GetInventory(id string) (*models.Inventory, error) {
	return s.InventoryModel.GetByID(id)
}

func (s *InventoryService) ListInventories(userID string) ([]*models.Inventory, error) {
	return s.InventoryModel.ListByUserID(userID)
}
