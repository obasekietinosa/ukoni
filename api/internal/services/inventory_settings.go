package services

import (
	"context"
	"database/sql"
	"errors"
	"ukoni/internal/database"
	"ukoni/internal/models"
)

type InventorySettingsService struct {
	InventorySettingsModel *models.InventorySettingsModel
	InventoryModel         *models.InventoryModel
	MembershipModel        *models.MembershipModel
	DB                     database.DBTX
}

func (s *InventorySettingsService) checkAdminOrOwner(ctx context.Context, userID, inventoryID string) (bool, error) {
	// Check owner
	inv, err := s.InventoryModel.GetByID(inventoryID)
	if err != nil {
		return false, err
	}
	if inv.OwnerUserID == userID {
		return true, nil
	}

	// Check member
	member, err := s.MembershipModel.GetMembership(inventoryID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	if member.Role == "admin" {
		return true, nil
	}

	return false, nil
}

func (s *InventorySettingsService) GetSettings(ctx context.Context, userID, inventoryID string) (*models.InventorySettings, error) {
	isAdmin, err := s.checkAdminOrOwner(ctx, userID, inventoryID)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, errors.New("unauthorized")
	}

	settings, err := s.InventorySettingsModel.Get(ctx, inventoryID)
	if err != nil {
		return nil, err
	}

	if settings == nil {
		return &models.InventorySettings{
			InventoryID: inventoryID,
		}, nil
	}

	return settings, nil
}

func (s *InventorySettingsService) UpdateSettings(ctx context.Context, userID, inventoryID string, settings *models.InventorySettings) (*models.InventorySettings, error) {
	isAdmin, err := s.checkAdminOrOwner(ctx, userID, inventoryID)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, errors.New("unauthorized")
	}

	settings.InventoryID = inventoryID
	err = s.InventorySettingsModel.Upsert(ctx, s.DB, settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}
