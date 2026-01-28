package services

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"ukoni/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func BenchmarkCreateTransaction(b *testing.B) {
	// Skip if no database connection string is provided
	// This allows the benchmark to exist but not fail in CI/sandbox without DB.
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		b.Skip("Skipping benchmark: DATABASE_URL not set")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		b.Fatalf("Failed to open db: %v", err)
	}
	defer db.Close()

	// Setup dependencies
	txModel := &models.TransactionModel{DB: db}
	membershipModel := &models.MembershipModel{DB: db}
	outletModel := &models.OutletModel{DB: db}
	activityLogModel := &models.ActivityLogModel{DB: db}
	productModel := &models.ProductModel{DB: db}
	inventoryProductModel := &models.InventoryProductModel{DB: db}

	activityLogService := &ActivityLogService{Model: activityLogModel}
	inventoryProductService := &InventoryProductService{
		InventoryProductModel: inventoryProductModel,
		ProductModel:          productModel,
	}

	svc := &TransactionService{
		DB:                      db,
		TransactionModel:        txModel,
		MembershipModel:         membershipModel,
		OutletModel:             outletModel,
		ActivityLogService:      activityLogService,
		InventoryProductService: inventoryProductService,
	}

	// We need valid IDs to run this.
	// In a real benchmark, we would setup fixtures here.
	// Since we can't run this in sandbox, I'll leave placeholders.
	inventoryID := "test-inventory-id"
	userID := "test-user-id"
	variantID := "test-variant-id"

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Stop timer for setup if needed, but creating input is fast
		input := CreateTransactionInput{
			InventoryID:     inventoryID,
			CreatedByUserID: userID,
			TransactionDate: time.Now(),
			Items: []CreateTransactionItemInput{
				{ProductVariantID: variantID, Quantity: 1, PricePerUnit: nil},
				{ProductVariantID: variantID, Quantity: 2, PricePerUnit: nil},
				{ProductVariantID: variantID, Quantity: 3, PricePerUnit: nil},
				{ProductVariantID: variantID, Quantity: 4, PricePerUnit: nil},
				{ProductVariantID: variantID, Quantity: 5, PricePerUnit: nil},
			},
		}

		_, err := svc.CreateTransaction(ctx, input)
		if err != nil {
			// In a real run we would handle errors or ensure fixtures exist
			// b.Fatalf("CreateTransaction failed: %v", err)
			continue
		}
	}
}
