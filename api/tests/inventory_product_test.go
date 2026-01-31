package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"ukoni/internal/models"
	"ukoni/internal/services"

	"github.com/golang-jwt/jwt/v5"
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
	canonical := &models.CanonicalProduct{
		Name:        "Milk",
		InventoryID: inventory.ID,
	}
	err = cpModel.Create(ctx, testDB, canonical)
	require.NoError(t, err)

	product := &models.Product{
		CanonicalProductID: &canonical.ID,
		Name:               "Milk Brand X",
		InventoryID:        inventory.ID,
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
		ID:          "tx-123",
		InventoryID: inventory.ID,
	}

	items := []*models.TransactionItem{
		{
			ProductVariantID: variant.ID,
			Quantity:         2.0,
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

func TestInventoryProduct_ListWithDetails(t *testing.T) {
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

	// 1. Setup Data
	user := &models.User{
		Email:        "test2@example.com",
		Name:         "Test User 2",
		PasswordHash: "password",
	}
	require.NoError(t, userModel.Insert(user))

	inventory := &models.Inventory{Name: "List Inventory", OwnerUserID: user.ID}
	require.NoError(t, inventoryModel.Create(ctx, testDB, inventory))

	canonical := &models.CanonicalProduct{Name: "Rice", InventoryID: inventory.ID}
	require.NoError(t, cpModel.Create(ctx, testDB, canonical))

	tilda := "Tilda"
	product := &models.Product{CanonicalProductID: &canonical.ID, Name: "Tilda Rice", Brand: &tilda, InventoryID: inventory.ID}
	require.NoError(t, productModel.Create(ctx, testDB, product))

	size := 1.0
	unit := "kg"
	variant := &models.ProductVariant{ProductID: product.ID, VariantName: "1kg Bag", Unit: &unit, Size: &size}
	require.NoError(t, productModel.CreateVariant(ctx, testDB, variant))

	// 2. Insert Inventory Product manually
	err := inventoryProductModel.Upsert(ctx, testDB, inventory.ID, variant.ID, 5.0, &unit) // 5kg
	require.NoError(t, err)

	// 3. Test ListWithDetails
	details, err := inventoryProductModel.ListWithDetails(ctx, inventory.ID, 10, 0)
	require.NoError(t, err)
	require.Len(t, details, 1)

	d := details[0]
	assert.Equal(t, canonical.Name, d.CanonicalProductName)
	assert.Equal(t, "Tilda", *d.BrandName)
	assert.Equal(t, "1kg Bag", d.VariantName)
	assert.Equal(t, 5.0, d.Quantity)
	assert.Equal(t, "kg", *d.Unit)
}

func TestInventoryProduct_ListEndpoint(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: no database connection")
	}
	clearDB()
	ctx := context.Background()

	// Setup data via models directly
	userModel := &models.UserModel{DB: testDB}
	inventoryModel := &models.InventoryModel{DB: testDB}
	productModel := &models.ProductModel{DB: testDB}
	inventoryProductModel := &models.InventoryProductModel{DB: testDB}
	cpModel := &models.CanonicalProductModel{DB: testDB}
	memModel := &models.MembershipModel{DB: testDB}

	user := &models.User{Email: "user@test.com", Name: "User", PasswordHash: "pass"}
	require.NoError(t, userModel.Insert(user))
	inventory := &models.Inventory{Name: "Inv", OwnerUserID: user.ID}
	require.NoError(t, inventoryModel.Create(ctx, testDB, inventory))
	// Add membership
	require.NoError(t, memModel.AddMember(ctx, testDB, inventory.ID, user.ID, "admin"))

	// Create Product
	canonical := &models.CanonicalProduct{Name: "Beans", InventoryID: inventory.ID}
	require.NoError(t, cpModel.Create(ctx, testDB, canonical))

	heinz := "Heinz"
	product := &models.Product{CanonicalProductID: &canonical.ID, Name: "Heinz Beans", Brand: &heinz, InventoryID: inventory.ID}
	require.NoError(t, productModel.Create(ctx, testDB, product))
	unit := "can"
	size := 1.0
	variant := &models.ProductVariant{ProductID: product.ID, VariantName: "400g", Unit: &unit, Size: &size}
	require.NoError(t, productModel.CreateVariant(ctx, testDB, variant))
	require.NoError(t, inventoryProductModel.Upsert(ctx, testDB, inventory.ID, variant.ID, 10.0, &unit))

	// Test Request
	router := setupRouter()

	req := httptest.NewRequest("GET", "/inventories/"+inventory.ID+"/inventory-products", nil)

	// Create token
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	token, err := tokenObj.SignedString([]byte(cfg.JWTSecret))
	require.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var response []*models.InventoryProductDetail
	err = json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	require.Len(t, response, 1)
	assert.Equal(t, "Beans", response[0].CanonicalProductName)
	assert.Equal(t, 10.0, response[0].Quantity)
}
