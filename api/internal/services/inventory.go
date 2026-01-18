package services

import (
	"errors"
	"ukoni/internal/models"
)

type InventoryService struct {
	InventoryModel *models.InventoryModel
}

func (s *InventoryService) CreateInventory(userID, name string) (*models.Inventory, error) {
	if name == "" {
		return nil, errors.New("inventory name cannot be empty")
	}

	inventory := &models.Inventory{
		Name:        name,
		OwnerUserID: userID,
	}

	if err := s.InventoryModel.Create(inventory); err != nil {
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
