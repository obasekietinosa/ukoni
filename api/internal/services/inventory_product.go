package services

import (
	"context"
	"fmt"
	"ukoni/internal/database"
	"ukoni/internal/models"
)

type InventoryProductService struct {
	InventoryProductModel *models.InventoryProductModel
	ProductModel          *models.ProductModel
}

func (s *InventoryProductService) UpdateFromTransaction(ctx context.Context, dbtx database.DBTX, transaction *models.Transaction, items []*models.TransactionItem) error {
	for _, item := range items {
		// Fetch variant
		variant, err := s.ProductModel.GetVariant(ctx, item.ProductVariantID)
		if err != nil {
			return fmt.Errorf("failed to get variant %s: %w", item.ProductVariantID, err)
		}
		if variant == nil {
			// Skip unknown variants, or should this be an error?
			// For data integrity, it should probably be an error, but if it was soft deleted during transaction?
			// Let's error.
			return fmt.Errorf("variant %s not found", item.ProductVariantID)
		}

		// Calculate quantity to add
		// Inventory Quantity = Item Quantity * (Variant Size if present else 1)
		qtyChange := item.Quantity
		if variant.Size != nil {
			qtyChange = item.Quantity * (*variant.Size)
		}

		// Update inventory
		err = s.InventoryProductModel.Upsert(ctx, dbtx, transaction.InventoryID, item.ProductVariantID, qtyChange, variant.Unit)
		if err != nil {
			return fmt.Errorf("failed to upsert inventory product: %w", err)
		}
	}
	return nil
}
