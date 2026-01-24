package services

import (
	"context"
	"ukoni/internal/database"
	"ukoni/internal/models"
)

type ActivityLogService struct {
	Model *models.ActivityLogModel
}

func (s *ActivityLogService) LogActivity(ctx context.Context, dbtx database.DBTX, inventoryID, userID *string, action, entityType string, entityID *string, metadata map[string]interface{}) error {
	logEntry := &models.ActivityLog{
		InventoryID: inventoryID,
		UserID:      userID,
		Action:      action,
		EntityType:  entityType,
		EntityID:    entityID,
		Metadata:    metadata,
	}

	return s.Model.Create(ctx, dbtx, logEntry)
}
