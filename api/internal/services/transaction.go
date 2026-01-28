package services

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"ukoni/internal/models"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
)

type TransactionService struct {
	DB                      *sql.DB
	TransactionModel        *models.TransactionModel
	MembershipModel         *models.MembershipModel
	OutletModel             *models.OutletModel
	ActivityLogService      *ActivityLogService
	InventoryProductService *InventoryProductService
}

type CreateTransactionInput struct {
	InventoryID     string
	OutletID        *string
	CreatedByUserID string
	TransactionDate time.Time
	Items           []CreateTransactionItemInput
}

type CreateTransactionItemInput struct {
	ProductVariantID   string
	Quantity           float64
	PricePerUnit       *float64
	ShoppingListItemID *string
}

type TransactionWithItems struct {
	*models.Transaction
	Items []*models.TransactionItem `json:"items"`
}

func (s *TransactionService) CreateTransaction(ctx context.Context, input CreateTransactionInput) (*models.Transaction, error) {
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

	// Validate outlet if present
	if input.OutletID != nil {
		_, err := s.OutletModel.Get(*input.OutletID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("outlet not found")
			}
			return nil, err
		}
	}

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Calculate total amount
	var totalAmount float64
	for _, item := range input.Items {
		if item.PricePerUnit != nil {
			totalAmount += *item.PricePerUnit * item.Quantity
		}
	}

	t := &models.Transaction{
		InventoryID:     input.InventoryID,
		OutletID:        input.OutletID,
		CreatedByUserID: input.CreatedByUserID,
		TransactionDate: input.TransactionDate,
		TotalAmount:     &totalAmount,
	}

	if err := s.TransactionModel.Create(ctx, tx, t); err != nil {
		return nil, err
	}

	createdItems := make([]*models.TransactionItem, 0, len(input.Items))
	for _, itemInput := range input.Items {
		item := &models.TransactionItem{
			TransactionID:      t.ID,
			ProductVariantID:   itemInput.ProductVariantID,
			Quantity:           itemInput.Quantity,
			PricePerUnit:       itemInput.PricePerUnit,
			ShoppingListItemID: itemInput.ShoppingListItemID,
		}
		createdItems = append(createdItems, item)
	}

	if err := s.TransactionModel.CreateItems(ctx, tx, createdItems); err != nil {
		return nil, err
	}

	// Update inventory
	if err := s.InventoryProductService.UpdateFromTransaction(ctx, tx, t, createdItems); err != nil {
		return nil, err
	}

	// Log activity
	// We use separate context/transaction for logging if we want it to persist even if main tx fails?
	// But here we log success. If commit fails, we don't log.
	// But ActivityLogService.LogActivity takes DBTX. If we use 'tx', it commits with it.
	// But we already committed 'tx'.
	// So we should log BEFORE commit?
	// Or use s.DB.
	// Typically we want the log to be part of the atomic transaction.
	if err := s.ActivityLogService.LogActivity(ctx, tx, &input.InventoryID, &input.CreatedByUserID, "transaction.created", "transaction", &t.ID, map[string]interface{}{
		"item_count":   len(input.Items),
		"total_amount": totalAmount,
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return t, nil
}

func (s *TransactionService) ListTransactions(ctx context.Context, inventoryID, userID string, limit, offset int) ([]*models.Transaction, error) {
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

	return s.TransactionModel.ListByInventory(ctx, inventoryID, limit, offset)
}

func (s *TransactionService) GetTransaction(ctx context.Context, transactionID, userID string) (*TransactionWithItems, error) {
	t, err := s.TransactionModel.GetByID(ctx, transactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}

	// Validate membership
	member, err := s.MembershipModel.GetMembership(t.InventoryID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user is not a member of this inventory")
		}
		return nil, err
	}
	if member == nil {
		return nil, errors.New("user is not a member of this inventory")
	}

	items, err := s.TransactionModel.GetItems(ctx, transactionID)
	if err != nil {
		return nil, err
	}

	return &TransactionWithItems{
		Transaction: t,
		Items:       items,
	}, nil
}
