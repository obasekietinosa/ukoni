package tests

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"ukoni/internal/models"
	"ukoni/internal/services"

	"github.com/google/uuid"
)

func seedPlanBenchmarkData(b *testing.B, s *services.PlanService, userID string) string {
	ctx := context.Background()

	// Create inventory
	inventoryID := uuid.New().String()
	_, err := testDB.Exec("INSERT INTO inventories (id, name, owner_user_id) VALUES ($1, $2, $3)", inventoryID, "Bench Plan Inventory", userID)
	if err != nil {
		b.Fatalf("Failed to create inventory: %v", err)
	}

	// Create Membership
	_, err = testDB.Exec("INSERT INTO inventory_memberships (inventory_id, user_id, role) VALUES ($1, $2, $3)", inventoryID, userID, "admin")
	if err != nil {
		b.Fatalf("Failed to create membership: %v", err)
	}

	// Create canonical product
	cpID := uuid.New().String()
	_, err = testDB.Exec("INSERT INTO canonical_products (id, name, inventory_id) VALUES ($1, $2, $3)", cpID, "Bench CP", inventoryID)
	if err != nil {
		b.Fatalf("Failed to create canonical product: %v", err)
	}

	// Create product
	pID := uuid.New().String()
	_, err = testDB.Exec("INSERT INTO products (id, canonical_product_id, name, inventory_id) VALUES ($1, $2, $3, $4)", pID, cpID, "Bench Product", inventoryID)
	if err != nil {
		b.Fatalf("Failed to create product: %v", err)
	}

	// Create variant
	vID := uuid.New().String()
	_, err = testDB.Exec("INSERT INTO product_variants (id, product_id, variant_name) VALUES ($1, $2, $3)", vID, pID, "Bench Variant")
	if err != nil {
		b.Fatalf("Failed to create variant: %v", err)
	}

	// Create Plan
	plan, err := s.CreatePlan(ctx, inventoryID, nil, "Bench Plan", "Description")
	if err != nil {
		b.Fatalf("Failed to create plan: %v", err)
	}

	// Add 20 items
	targets := []struct {
		Type string
		ID   string
	}{
		{"canonical_product", cpID},
		{"product", pID},
		{"product_variant", vID},
	}

	for i := 0; i < 20; i++ {
		target := targets[rand.Intn(len(targets))]
		qty := 1.0
		_, err := s.AddItem(ctx, plan.ID, target.Type, target.ID, &qty, "kg", "note")
		if err != nil {
			b.Fatalf("Failed to add item: %v", err)
		}
	}

	// Add 5 child plans
	for i := 0; i < 5; i++ {
		_, err := s.CreatePlan(ctx, inventoryID, &plan.ID, fmt.Sprintf("Child Plan %d", i), "")
		if err != nil {
			b.Fatalf("Failed to create child plan: %v", err)
		}
	}

	// Link shopping list
	sl := &models.ShoppingList{
		InventoryID: inventoryID,
		Name:        "Bench List",
		CreatedBy:   userID,
	}
	err = s.ShoppingListModel.CreateList(ctx, sl)
	if err != nil {
		b.Fatalf("Failed to create shopping list: %v", err)
	}

	err = s.LinkShoppingList(ctx, plan.ID, sl.ID)
	if err != nil {
		b.Fatalf("Failed to link shopping list: %v", err)
	}

	return plan.ID
}

func BenchmarkGetPlan(b *testing.B) {
	if testDB == nil {
		b.Skip("Database not initialized")
	}
	clearDB()

	// Setup service
	planModel := &models.PlanModel{DB: testDB}
	invModel := &models.InventoryModel{DB: testDB}
	prodModel := &models.ProductModel{DB: testDB}
	cpModel := &models.CanonicalProductModel{DB: testDB}
	slModel := &models.ShoppingListModel{DB: testDB}

	s := &services.PlanService{
		DB:                    testDB,
		PlanModel:             planModel,
		InventoryModel:        invModel,
		ProductModel:          prodModel,
		CanonicalProductModel: cpModel,
		ShoppingListModel:     slModel,
	}

	// Create user
	userID := uuid.New().String()
	_, err := testDB.Exec("INSERT INTO users (id, email, name, password_hash) VALUES ($1, $2, $3, $4)", userID, "benchplan@example.com", "Bench Plan User", "hash")
	if err != nil {
		b.Fatalf("Failed to create user: %v", err)
	}

	planID := seedPlanBenchmarkData(b, s, userID)

	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := s.GetPlan(ctx, planID)
		if err != nil {
			b.Fatalf("GetPlan failed: %v", err)
		}
	}
}

func BenchmarkGetPlanSummary(b *testing.B) {
	if testDB == nil {
		b.Skip("Database not initialized")
	}
	clearDB()

	// Setup service
	planModel := &models.PlanModel{DB: testDB}
	invModel := &models.InventoryModel{DB: testDB}
	prodModel := &models.ProductModel{DB: testDB}
	cpModel := &models.CanonicalProductModel{DB: testDB}
	slModel := &models.ShoppingListModel{DB: testDB}

	s := &services.PlanService{
		DB:                    testDB,
		PlanModel:             planModel,
		InventoryModel:        invModel,
		ProductModel:          prodModel,
		CanonicalProductModel: cpModel,
		ShoppingListModel:     slModel,
	}

	// Create user
	userID := uuid.New().String()
	_, err := testDB.Exec("INSERT INTO users (id, email, name, password_hash) VALUES ($1, $2, $3, $4)", userID, "benchplansummary@example.com", "Bench Plan Summary User", "hash")
	if err != nil {
		b.Fatalf("Failed to create user: %v", err)
	}

	planID := seedPlanBenchmarkData(b, s, userID)

	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := s.GetPlanSummary(ctx, planID)
		if err != nil {
			b.Fatalf("GetPlanSummary failed: %v", err)
		}
	}
}
