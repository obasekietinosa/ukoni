package tests

import (
	"context"
	"testing"
	"ukoni/internal/models"
	"ukoni/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInventoryProduct_UpdateFromTransaction(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: no database connection")
	}
	clearDB()
	ctx := context.Background()

	// Setup dependencies
	userModel := &models.UserModel{DB: testDB}
	inventoryModel := &models.InventoryModel{DB: testDB}
	productModel := &models.ProductModel{DB: testDB}
	inventoryProductModel := &models.InventoryProductModel{DB: testDB}
	cpModel := &models.CanonicalProductModel{DB: testDB}

	svc := &services.InventoryProductService{
		InventoryProductModel: inventoryProductModel,
		ProductModel:          productModel,
	}

	// 1. Create User
	user := &models.User{
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: "password",
	}
	err := userModel.Insert(user)
	require.NoError(t, err)

	// 2. Create Inventory
	inventory := &models.Inventory{
		Name:        "Test Inventory",
		OwnerUserID: user.ID,
	}
	err = inventoryModel.Create(ctx, testDB, inventory)
	require.NoError(t, err)

	// 3. Create Product & Variant
	canonical := &models.CanonicalProduct{Name: "Milk"}
	err = cpModel.Create(ctx, testDB, canonical)
	require.NoError(t, err)

	product := &models.Product{
		CanonicalProductID: &canonical.ID,
		Name:               "Milk Brand X",
	}
	err = productModel.Create(ctx, testDB, product)
	require.NoError(t, err)

	size := 1.5
	unit := "L"
	variant := &models.ProductVariant{
		ProductID:   product.ID,
		VariantName: "1.5L Bottle",
		Unit:        &unit,
		Size:        &size,
	}
	err = productModel.CreateVariant(ctx, testDB, variant)
	require.NoError(t, err)

	// 4. Create Transaction Item (Mocking transaction flow)
	tx := &models.Transaction{
		ID:          "tx-123", // The ID is used for logging but not FK in this service call, unless needed?
        // Service doesn't use tx.ID for logic, only passed it.
        // But UpdateFromTransaction passes 'tx' to... nowhere?
        // Ah, UpdateFromTransaction uses transaction.InventoryID.
		InventoryID: inventory.ID,
	}

	items := []*models.TransactionItem{
		{
			ProductVariantID: variant.ID,
			Quantity:         2.0, // Buying 2 bottles
		},
	}

	// 5. Call Service
	err = svc.UpdateFromTransaction(ctx, testDB, tx, items)
	require.NoError(t, err)

	// 6. Verify Inventory
	ip, err := inventoryProductModel.Get(ctx, inventory.ID, variant.ID)
	require.NoError(t, err)
	require.NotNil(t, ip)

	// Quantity should be 2.0 * 1.5 = 3.0
	assert.Equal(t, 3.0, ip.Quantity)
	assert.Equal(t, "L", *ip.Unit)

	// 7. Update again (buy 1 more)
	items2 := []*models.TransactionItem{
		{
			ProductVariantID: variant.ID,
			Quantity:         1.0,
		},
	}
	err = svc.UpdateFromTransaction(ctx, testDB, tx, items2)
	require.NoError(t, err)

	ip, err = inventoryProductModel.Get(ctx, inventory.ID, variant.ID)
	require.NoError(t, err)
	assert.Equal(t, 4.5, ip.Quantity) // 3.0 + 1.5
}
