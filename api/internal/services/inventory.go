package services

import (
	"errors"
	"ukoni/internal/models"
)

type InventoryService struct {
	InventoryModel  *models.InventoryModel
	MembershipModel *models.MembershipModel
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

	// Add owner as an admin member
	err := s.MembershipModel.AddMember(inventory.ID, userID, "admin")
	if err != nil {
		// potential consistency issue here if this fails but inventory was created.
		// In a real app we'd use a transaction spanning both models.
		// For now, we'll log/return error.
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
