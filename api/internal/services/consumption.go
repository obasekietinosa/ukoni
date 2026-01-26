package services

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"ukoni/internal/models"
)

type ConsumptionService struct {
	DB                 *sql.DB
	ConsumptionModel   *models.ConsumptionModel
	MembershipModel    *models.MembershipModel
	ActivityLogService *ActivityLogService
}

type CreateConsumptionInput struct {
	InventoryID        string
	CanonicalProductID *string
	CreatedByUserID    string
	Quantity           *float64
	Unit               *string
	Note               *string
	Source             string
	ConsumedAt         time.Time
}

func (s *ConsumptionService) CreateConsumption(ctx context.Context, input CreateConsumptionInput) (*models.ConsumptionEvent, error) {
	// Validate membership
	member, err := s.MembershipModel.GetMembership(input.InventoryID, input.CreatedByUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user is not a member of this inventory")
		}
		return nil, err
	}
	if member == nil {
		return nil, errors.New("user is not a member of this inventory")
	}

	event := &models.ConsumptionEvent{
		InventoryID:        input.InventoryID,
		CanonicalProductID: input.CanonicalProductID,
		CreatedByUserID:    &input.CreatedByUserID,
		Quantity:           input.Quantity,
		Unit:               input.Unit,
		Note:               input.Note,
		Source:             input.Source,
		ConsumedAt:         input.ConsumedAt,
	}

	if event.Source == "" {
		event.Source = "manual"
	}

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := s.ConsumptionModel.Create(ctx, tx, event); err != nil {
		return nil, err
	}

	if err := s.ActivityLogService.LogActivity(ctx, tx, &input.InventoryID, &input.CreatedByUserID, "consumption.created", "consumption_event", &event.ID, map[string]interface{}{
		"canonical_product_id": input.CanonicalProductID,
		"quantity":             input.Quantity,
		"unit":                 input.Unit,
		"source":               input.Source,
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return event, nil
}

func (s *ConsumptionService) ListConsumptionEvents(ctx context.Context, inventoryID, userID string, limit, offset int) ([]*models.ConsumptionEvent, error) {
	// Validate membership
	member, err := s.MembershipModel.GetMembership(inventoryID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user is not a member of this inventory")
		}
		return nil, err
	}
	if member == nil {
		return nil, errors.New("user is not a member of this inventory")
	}

	return s.ConsumptionModel.List(ctx, inventoryID, limit, offset)
}
